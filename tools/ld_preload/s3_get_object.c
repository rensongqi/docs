#include <libs3.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
// Callback function to handle the response
static S3Status responsePropertiesCallback(
    const S3ResponseProperties *properties,
    void *callbackData)
{
    return S3StatusOK;
}
// Callback function to handle errors
static void responseCompleteCallback(
    S3Status status,
    const S3ErrorDetails *error,
    void *callbackData)
{
    if (status == S3StatusOK)
    {
        printf("Object downloaded successfully.\n");
    }
    else
    {
        printf("Failed to get object. Error: %s\n", S3_get_status_name(status));
        if (error && error->message)
        {
            printf("Error message: %s\n", error->message);
        }
    }
}
// Callback function to receive object data
static S3Status getObjectDataCallback(
    int bufferSize,
    const char *buffer,
    void *callbackData)
{
    FILE *outfile = (FILE *)callbackData;
    size_t wrote = fwrite(buffer, 1, bufferSize, outfile);
    return ((wrote < (size_t)bufferSize) ? S3StatusAbortedByCallback : S3StatusOK);
}
int main(int argc, char **argv)
{
    const char *host = "172.16.1.10:9000";
    const char *accessKeyId = "SMNmxxxxxxxxxxxxx";
    const char *secretAccessKey = "oKWoyEctxxxxxxxxxxxxxxxxxxx";
    const char *bucketName = "devops";
    const char *objectKey = "file.special";
    const char *region = "xxxxxx";

    char my_buffer[1024];
    uint64_t startByte = 0;
    uint64_t byteCount = 0;
    FILE *outfile = tmpfile();
    // Initialize libs3
    S3_initialize("s3", S3_INIT_ALL, region);
    // Set up the bucket context
    S3BucketContext bucketContext = {
        .hostName = host,
        .bucketName = bucketName,
        .protocol = S3ProtocolHTTP,
        .authRegion = region,
        .accessKeyId = accessKeyId,
        .secretAccessKey = secretAccessKey,
        .uriStyle = S3UriStylePath};
    // Set up the get object handler
    S3GetObjectHandler getObjectHandler = {
        .responseHandler = {
            .propertiesCallback = responsePropertiesCallback,
            .completeCallback = responseCompleteCallback},
        .getObjectDataCallback = getObjectDataCallback};
    // Get the object (no need to assign a return value, since S3_get_object returns void)
    S3_get_object(&bucketContext, objectKey, NULL, startByte,
                  byteCount, NULL, 0, &getObjectHandler, outfile);
    // Reset file pointer to beginning of file
    rewind(outfile);
    // Read and print the content
    printf("Object content:\n");
    while (fgets(my_buffer, sizeof(my_buffer), outfile) != NULL)
    {
        printf("%s", my_buffer);
    }
    // Clean up
    fclose(outfile);
    S3_deinitialize();
    return 0;
}