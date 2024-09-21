#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>
#include <sys/xattr.h>
#include <unistd.h>
#include <fcntl.h>
#include <libs3.h>

typedef char *(*fgets_type)(char *, int, FILE *);
typedef size_t (*fread_type)(void *, size_t, size_t, FILE *);
typedef int (*fscanf_type)(FILE *, const char *, ...);
typedef int (*orig_open_f_type) (const char *pathname, int flags);
typedef ssize_t (*read_type)(int fd, void *buf, size_t count);

// 用于跟踪文件状态的结构
typedef struct {
    FILE *redirected_file;
    int original_fd;
    int eof_reached;
    int redirected_fd;
} FileStatus;
#define MAX_FILES 1024
FileStatus file_status[MAX_FILES] = {0};

const char *host = "172.16.1.10:9000";
const char *accessKeyId = "SMNmxxxxxxxxxxxxx";
const char *secretAccessKey = "oKWoyEctxxxxxxxxxxxxxxxxxxx";
const char *bucketName = "devops";
const char *region = "xxxxx";

// S3 回调函数
static S3Status responsePropertiesCallback(const S3ResponseProperties *properties, void *callbackData) {
    return S3StatusOK;
}
static void responseCompleteCallback(S3Status status, const S3ErrorDetails *error, void *callbackData) {
    if (status != S3StatusOK) {
        fprintf(stderr, "Failed to get object. Error: %s\n", S3_get_status_name(status));
        if (error && error->message) {
            fprintf(stderr, "Error message: %s\n", error->message);
        }
    }
}
static S3Status getObjectDataCallback(int bufferSize, const char *buffer, void *callbackData) {
    FILE *outfile = (FILE *)callbackData;
    size_t wrote = fwrite(buffer, 1, bufferSize, outfile);
    return (wrote < (size_t)bufferSize) ? S3StatusAbortedByCallback : S3StatusOK;
}

// 从 S3 下载对象并保存到临时文件
int fetch_object_from_s3(const char *objectKey, const char *bucketName, FILE *outfile) {
    S3_initialize("s3", S3_INIT_ALL, region);
    S3BucketContext bucketContext = {
        .hostName = host,
        .bucketName = bucketName,
        .protocol = S3ProtocolHTTP,
        .authRegion = region,
        .accessKeyId = accessKeyId,
        .secretAccessKey = secretAccessKey,
        .uriStyle = S3UriStylePath
    };
    S3GetObjectHandler getObjectHandler = {
        .responseHandler = {
            .propertiesCallback = responsePropertiesCallback,
            .completeCallback = responseCompleteCallback
        },
        .getObjectDataCallback = getObjectDataCallback
    };
    uint64_t startByte = 0;
    uint64_t byteCount = 0;
    S3_get_object(&bucketContext, objectKey, NULL, startByte, byteCount, NULL, 0, &getObjectHandler, outfile);
    S3_deinitialize();
    return 0;
}

// 拦截fgets的请求
char *fgets(char *str, int n, FILE *stream) {
    static fgets_type orig_fgets = NULL;
    if (!orig_fgets) {
        orig_fgets = (fgets_type)dlsym(RTLD_NEXT, "fgets");
        if (!orig_fgets) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }

    int fd = fileno(stream); // 获取文件描述符

    // 检查文件描述符是否有效
    if (fd < 0 || fd >= MAX_FILES) {
        return orig_fgets(str, n, stream); // 如果文件描述符无效，继续原始操作
    }

    // 如果这是一个新的文件描述符，检查是否需要重定向
    if (!file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        char bucketName[256];
        char objectKey[256];
        ssize_t s3_bucket = fgetxattr(fd, "user.s3_bucket", bucketName, sizeof(bucketName) - 1);
        ssize_t s3_object = fgetxattr(fd, "user.s3_object", objectKey, sizeof(objectKey) - 1);
        if (s3_bucket > 0 && s3_object > 0) {
            bucketName[s3_bucket] = '\0';
            objectKey[s3_object] = '\0';
            FILE *s3_file = tmpfile();    // 创建临时文件
            if (fetch_object_from_s3(objectKey, bucketName, s3_file) == 0) {
                // S3 对象获取成功，重定向读取
                file_status[fd].redirected_file = s3_file;
                rewind(s3_file);  // 确保临时文件的文件指针重置为开头
            } else {
                return orig_fgets(str, n, stream);  // 如果获取失败，继续原始文件操作
            }
        }
    }

    // 如果有重定向文件且未到达 EOF
    if (file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        // 从临时文件读取数据到缓冲区
        char *result = orig_fgets(str, n, file_status[fd].redirected_file);
        if (!result) {
            // 临时文件已经读取完毕，标记 EOF 并关闭临时文件
            file_status[fd].eof_reached = 1;
            fclose(file_status[fd].redirected_file);
            file_status[fd].redirected_file = NULL;
        }
        return result;
    }

    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    return orig_fgets(str, n, stream);
}

// fgets2 从本地文件读取并返回数据给客户端程序
char *fgets2(char *str, int n, FILE *stream) {
    static fgets_type orig_fgets = NULL;
    if (!orig_fgets) {
        orig_fgets = (fgets_type)dlsym(RTLD_NEXT, "fgets");
        if (!orig_fgets) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }
    int fd = fileno(stream);
    if (fd < 0 || fd >= MAX_FILES) {
        return orig_fgets(str, n, stream);
    }
    // 如果这是一个新的文件流
    if (!file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        char filename[256];
        ssize_t attr_size = fgetxattr(fd, "user.dest", filename, sizeof(filename) - 1);
        if (attr_size > 0) {
            filename[attr_size] = '\0';  // Null-terminate the string
            FILE *new_file = fopen(filename, "r");
            if (new_file) {
                file_status[fd].redirected_file = new_file;
                file_status[fd].original_fd = fd;
            }
        }
    }
    // 如果有重定向文件且未到达EOF
    if (file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        char *result = orig_fgets(str, n, file_status[fd].redirected_file);
        if (result == NULL) {
            // 重定向文件已读完，标记EOF并关闭
            file_status[fd].eof_reached = 1;
            fclose(file_status[fd].redirected_file);
            file_status[fd].redirected_file = NULL;
            // 继续读取原始文件
            // return orig_fgets(str, n, stream);
        }
        return result;
    }
    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    return orig_fgets(str, n, stream);
}

// 拦截fread方法，返回客户端程序
size_t fread(void *ptr, size_t size, size_t count, FILE *stream) {
    static fread_type orig_fread = NULL;
    if (!orig_fread) {
        orig_fread = (fread_type)dlsym(RTLD_NEXT, "fread");
        if (!orig_fread) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }
    int fd = fileno(stream);
    if (fd < 0 || fd >= MAX_FILES) {
        return orig_fread(ptr, size, count, stream);
    }
    // 如果这是一个新的文件流
    if (!file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        char filename[256];
        ssize_t attr_size = fgetxattr(fd, "user.dest", filename, sizeof(filename) - 1);
        if (attr_size > 0) {
            filename[attr_size] = '\0';  // Null-terminate the string
            FILE *new_file = fopen(filename, "r");
            if (new_file) {
                file_status[fd].redirected_file = new_file;
                file_status[fd].original_fd = fd;
            }
        }
    }
    // 如果有重定向文件且未到达EOF
    if (file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        size_t result = orig_fread(ptr, size, count, file_status[fd].redirected_file);
        if (result == 0) {
            // 重定向文件已读完，标记EOF并关闭
            file_status[fd].eof_reached = 1;
            fclose(file_status[fd].redirected_file);
            file_status[fd].redirected_file = NULL;
        }
        return result;
    }
    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    return orig_fread(ptr, size, count, stream);
}

// 拦截fscanf方法
int fscanf(FILE *stream, const char *format, ...) {
    static fscanf_type orig_fscanf = NULL;
    if (!orig_fscanf) {
        orig_fscanf = (fscanf_type)dlsym(RTLD_NEXT, "fscanf");
        if (!orig_fscanf) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }
    int fd = fileno(stream);
    if (fd < 0 || fd >= MAX_FILES) {
        va_list args;
        va_start(args, format);
        int result = vfprintf(stream, format, args);
        va_end(args);
        return result;
    }
    // 如果这是一个新的文件流
    if (!file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        char filename[256];
        ssize_t attr_size = fgetxattr(fd, "user.dest", filename, sizeof(filename) - 1);
        if (attr_size > 0) {
            filename[attr_size] = '\0';  // Null-terminate the string
            FILE *new_file = fopen(filename, "r");
            if (new_file) {
                file_status[fd].redirected_file = new_file;
                file_status[fd].original_fd = fd;
            }
        }
    }
    // 如果有重定向文件且未到达EOF
    if (file_status[fd].redirected_file && !file_status[fd].eof_reached) {
        va_list args;
        va_start(args, format);
        int result = vfscanf(file_status[fd].redirected_file, format, args);
        va_end(args);
        if (result == EOF) {
            // 重定向文件已读完，标记EOF并关闭
            file_status[fd].eof_reached = 1;
            fclose(file_status[fd].redirected_file);
            file_status[fd].redirected_file = NULL;
        }
        return result;
    }
    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    va_list args;
    va_start(args, format);
    int result = vfscanf(stream, format, args);
    va_end(args);
    return result;
}

// 拦截open方法，ls、cat、cp等命令都会先触发open方法
// int open(const char *pathname, int flags, ...)
// {
//     // remember to include stdio.h!
//     printf("open():%s\n", pathname);
//     /* Some evil injected code goes here. */
//     orig_open_f_type orig_open;
//     orig_open = (orig_open_f_type) dlsym(RTLD_NEXT, "open");
//     return orig_open(pathname, flags);
// }
// FILE *fopen(const char *path, const char *mode) {
//     printf("fopen():%s\n", path);
//     FILE* (*original_fopen) (const char*, const char*);
//     original_fopen = dlsym(RTLD_NEXT, "fopen");
//     return (*original_fopen)(path, mode);
// }


// ssize_t read(int fd, void *buf, size_t count) {
//     printf("read(): fd=%d, count=%zu\n", fd, count);
//     orig_read_f_type orig_read;
//     orig_read = (orig_read_f_type) dlsym(RTLD_NEXT, "read");
//     return orig_read(fd, buf, count);
// }

//  拦截read方法，python和rust底层都基于read方法实现
ssize_t read(int fd, void *buf, size_t count) {
    static read_type orig_read = NULL;
    if (!orig_read) {
        orig_read = (read_type)dlsym(RTLD_NEXT, "read");
        if (!orig_read) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }
    // 检查文件描述符是否有效
    if (fd < 0 || fd >= MAX_FILES) {
        return orig_read(fd, buf, count);
    }
    // 如果这是一个新的文件描述符，检查是否需要重定向
    if (!file_status[fd].redirected_fd && !file_status[fd].eof_reached) {
        char bucketName[256];
        char objectKey[256];
        ssize_t s3_bucket = fgetxattr(fd, "user.s3_bucket", bucketName, sizeof(bucketName) - 1);
        ssize_t s3_object = fgetxattr(fd, "user.s3_object", objectKey, sizeof(objectKey) - 1);
        if (s3_bucket > 0 && s3_object > 0) {
            bucketName[s3_bucket] = '\0';
            objectKey[s3_object] = '\0';
            FILE *s3_file = tmpfile();    // 创建临时文件
            if (fetch_object_from_s3(objectKey, bucketName, s3_file) == 0) {
                // S3 对象获取成功，重定向读取
                file_status[fd].redirected_fd = fileno(s3_file);
                file_status[fd].redirected_file = s3_file;
                file_status[fd].original_fd = fd;
                rewind(s3_file);  // 确保临时文件的文件指针重置为开头
            } else {
                return orig_read(fd, buf, count);  // 如果获取失败，继续读取原始文件
            }
        }
    }
    // 如果有重定向文件且未到达 EOF
    if (file_status[fd].redirected_fd && !file_status[fd].eof_reached) {
        // 从临时文件读取数据到缓冲区
        ssize_t result = fread(buf, 1, count, file_status[fd].redirected_file);
        if (result == 0) {
            // 临时文件已经读取完毕，标记 EOF 并关闭临时文件
            file_status[fd].eof_reached = 1;
            fclose(file_status[fd].redirected_file);
            file_status[fd].redirected_fd = -1;
        }
        return result;
    }
    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    return orig_read(fd, buf, count);
}

// 拦截read方法，从本地读取文件返回给客户端程序
ssize_t read2(int fd, void *buf, size_t count) {
    static read_type orig_read = NULL;
    if (!orig_read) {
        orig_read = (read_type)dlsym(RTLD_NEXT, "read");
        if (!orig_read) {
            fprintf(stderr, "Error: %s\n", dlerror());
            exit(1);
        }
    }
    if (fd < 0 || fd >= MAX_FILES) {
        return orig_read(fd, buf, count);
    }
    // 如果这是一个新的文件描述符
    if (!file_status[fd].redirected_fd && !file_status[fd].eof_reached) {
        char filename[256];
        ssize_t attr_size = fgetxattr(fd, "user.dest", filename, sizeof(filename) - 1);
        if (attr_size > 0) {
            filename[attr_size] = '\0';  // Null-terminate the string
            int new_fd = open(filename, O_RDONLY);
            if (new_fd != -1) {
                file_status[fd].redirected_fd = new_fd;
                file_status[fd].original_fd = fd;
            }
        }
    }
    // 如果有重定向文件且未到达EOF
    if (file_status[fd].redirected_fd && !file_status[fd].eof_reached) {
        ssize_t result = orig_read(file_status[fd].redirected_fd, buf, count);
        if (result == 0) {
            // 重定向文件已读完，标记EOF并关闭
            file_status[fd].eof_reached = 1;
            close(file_status[fd].redirected_fd);
            file_status[fd].redirected_fd = -1;
        }
        return result;
    }
    // 如果没有重定向或已经读完重定向文件，则读取原始文件
    return orig_read(fd, buf, count);
}
