#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <infiniband/verbs.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#define MSG_SIZE    1024000
#define SERVER_PORT 12345

struct rdma_connection_info {
    uint32_t qp_num;
    uint16_t lid;
    uint32_t psn;
};
// RDMA 资源结构体
struct rdma_resources {
    struct ibv_context *context;
    struct ibv_pd *pd;
    struct ibv_mr *mr;
    struct ibv_cq *cq;
    struct ibv_qp *qp;
    char *buf;
    int sock;
};

// 初始化 RDMA 设备
static struct rdma_resources *init_rdma_resources() {
    struct rdma_resources *res = calloc(1, sizeof(struct rdma_resources));
    if (!res) {
        return NULL;
    }
    // 获取 RDMA 设备列表
    int num_devices;
    struct ibv_device **dev_list = ibv_get_device_list(&num_devices);
    if (!dev_list) {
        fprintf(stderr, "Failed to get IB devices list\n");
        goto cleanup;
    }
    // 打开第一个设备
    res->context = ibv_open_device(dev_list[0]);
    if (!res->context) {
        fprintf(stderr, "Failed to open IB device\n");
        goto cleanup;
    }

    // 分配保护域
    res->pd = ibv_alloc_pd(res->context);
    if (!res->pd) {
        fprintf(stderr, "Failed to allocate PD\n");
        goto cleanup;
    }
    // 创建完成队列
    res->cq = ibv_create_cq(res->context, 10, NULL, NULL, 0);
    if (!res->cq) {
        fprintf(stderr, "Failed to create CQ\n");
        goto cleanup;
    }
    // 分配内存缓冲区
    res->buf = malloc(MSG_SIZE);
    if (!res->buf) {
        fprintf(stderr, "Failed to allocate memory buffer\n");
        goto cleanup;
    }
    // 注册内存区域
    res->mr = ibv_reg_mr(res->pd, res->buf, MSG_SIZE,
                         IBV_ACCESS_LOCAL_WRITE | 
                         IBV_ACCESS_REMOTE_WRITE | 
                         IBV_ACCESS_REMOTE_READ);
    if (!res->mr) {
        fprintf(stderr, "Failed to register MR\n");
        goto cleanup;
    }
    ibv_free_device_list(dev_list);
    return res;
cleanup:
    if (res->mr) ibv_dereg_mr(res->mr);
    if (res->buf) free(res->buf);
    if (res->cq) ibv_destroy_cq(res->cq);
    if (res->pd) ibv_dealloc_pd(res->pd);
    if (res->context) ibv_close_device(res->context);
    if (dev_list) ibv_free_device_list(dev_list);
    free(res);
    return NULL;
}
// 创建队列对
static int create_qp(struct rdma_resources *res) {
    struct ibv_qp_init_attr qp_init_attr = {
        .send_cq = res->cq,
        .recv_cq = res->cq,
        .cap = {
            .max_send_wr = 10,
            .max_recv_wr = 10,
            .max_send_sge = 1,
            .max_recv_sge = 1,
        },
        .qp_type = IBV_QPT_RC,
    };
    res->qp = ibv_create_qp(res->pd, &qp_init_attr);
    if (!res->qp) {
        fprintf(stderr, "Failed to create QP\n");
        return -1;
    }
    return 0;
}
// 修改队列对状态为 INIT
static int modify_qp_to_init(struct ibv_qp *qp) {
    struct ibv_qp_attr attr = {
        .qp_state = IBV_QPS_INIT,
        .pkey_index = 0,
        .port_num = 1,
        .qp_access_flags = IBV_ACCESS_REMOTE_READ |
                          IBV_ACCESS_REMOTE_WRITE |
                          IBV_ACCESS_REMOTE_ATOMIC,
    };
    return ibv_modify_qp(qp, &attr,
                        IBV_QP_STATE |
                        IBV_QP_PKEY_INDEX |
                        IBV_QP_PORT |
                        IBV_QP_ACCESS_FLAGS);
}
// 修改队列对状态为 RTR
static int modify_qp_to_rtr(struct ibv_qp *qp, uint32_t remote_qpn, uint16_t dlid, uint32_t psn) {
    struct ibv_qp_attr attr = {
        .qp_state = IBV_QPS_RTR,
        .path_mtu = IBV_MTU_1024,
        .dest_qp_num = remote_qpn,
        .rq_psn = psn,
        .max_dest_rd_atomic = 1,
        .min_rnr_timer = 12,
        .ah_attr = {
            .dlid = dlid,
            .sl = 0,
            .src_path_bits = 0,
            .port_num = 1,
        },
    };
    return ibv_modify_qp(qp, &attr,
                        IBV_QP_STATE |
                        IBV_QP_AV |
                        IBV_QP_PATH_MTU |
                        IBV_QP_DEST_QPN |
                        IBV_QP_RQ_PSN |
                        IBV_QP_MAX_DEST_RD_ATOMIC |
                        IBV_QP_MIN_RNR_TIMER);
}
// 修改队列对状态为 RTS
static int modify_qp_to_rts(struct ibv_qp *qp, uint32_t psn) {
    struct ibv_qp_attr attr = {
        .qp_state = IBV_QPS_RTS,
        .timeout = 14,
        .retry_cnt = 7,
        .rnr_retry = 7,
        .sq_psn = psn,
        .max_rd_atomic = 1,
    };
    return ibv_modify_qp(qp, &attr,
                        IBV_QP_STATE |
                        IBV_QP_TIMEOUT |
                        IBV_QP_RETRY_CNT |
                        IBV_QP_RNR_RETRY |
                        IBV_QP_SQ_PSN |
                        IBV_QP_MAX_QP_RD_ATOMIC);
}

// 服务端实现
int run_server() {
    struct rdma_resources *res = init_rdma_resources();
    if (!res) {
        return -1;
    }
    int listen_sock = socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in server_addr = {
        .sin_family = AF_INET,
        .sin_port = htons(SERVER_PORT),
        .sin_addr.s_addr = INADDR_ANY,
    };
    bind(listen_sock, (struct sockaddr *)&server_addr, sizeof(server_addr));
    listen(listen_sock, 1);
    res->sock = accept(listen_sock, NULL, NULL);
    if (create_qp(res) != 0) {
        goto cleanup;
    }
    struct rdma_connection_info local_info = {
        .qp_num = res->qp->qp_num,
        .lid = 23,
        .psn = rand() & 0xffffff,
    };
    struct rdma_connection_info remote_info;
    write(res->sock, &local_info, sizeof(local_info));
    read(res->sock, &remote_info, sizeof(remote_info));
    if (modify_qp_to_init(res->qp) ||
        modify_qp_to_rtr(res->qp, remote_info.qp_num, remote_info.lid, remote_info.psn) ||
        modify_qp_to_rts(res->qp, local_info.psn)) {
        fprintf(stderr, "Failed to modify QP state\n");
        goto cleanup;
    }
    int out_fd = open("received_file", O_WRONLY | O_CREAT | O_TRUNC, 0666);
    if (out_fd < 0) {
        fprintf(stderr, "Failed to open output file: %s\n", strerror(errno));
        goto cleanup;
    }
    while (1) {
        struct ibv_recv_wr wr = {
            .wr_id = 1,
            .sg_list = &(struct ibv_sge){
                .addr = (uint64_t)res->buf,
                .length = MSG_SIZE,
                .lkey = res->mr->lkey,
            },
            .num_sge = 1,
        };
        struct ibv_recv_wr *bad_wr;
        if (ibv_post_recv(res->qp, &wr, &bad_wr)) {
            fprintf(stderr, "Failed to post receive request\n");
            goto cleanup;
        }
        struct ibv_wc wc;
        while (ibv_poll_cq(res->cq, 1, &wc) == 0);
        if (wc.status != IBV_WC_SUCCESS) {
            fprintf(stderr, "Receive failed with status %d\n", wc.status);
            goto cleanup;
        }
        write(out_fd, res->buf, wc.byte_len);
        if (wc.byte_len < MSG_SIZE) break; // End of file
    }
    close(out_fd);
cleanup:
    if (res->qp) ibv_destroy_qp(res->qp);
    if (res->mr) ibv_dereg_mr(res->mr);
    if (res->buf) free(res->buf);
    if (res->cq) ibv_destroy_cq(res->cq);
    if (res->pd) ibv_dealloc_pd(res->pd);
    if (res->context) ibv_close_device(res->context);
    close(res->sock);
    close(listen_sock);
    free(res);
    return 0;
}
// 客户端实现
int run_client(const char *server_ip, const char *file_path) {
    struct rdma_resources *res = init_rdma_resources();
    if (!res) {
        return -1;
    }
    res->sock = socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in server_addr = {
        .sin_family = AF_INET,
        .sin_port = htons(SERVER_PORT),
    };
    inet_pton(AF_INET, server_ip, &server_addr.sin_addr);
    if (connect(res->sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) != 0) {
        fprintf(stderr, "Failed to connect to server\n");
        goto cleanup;
    }
    if (create_qp(res) != 0) {
        goto cleanup;
    }
    struct rdma_connection_info local_info = {
        .qp_num = res->qp->qp_num,
        .lid = 11,
        .psn = rand() & 0xffffff,
    };
    struct rdma_connection_info remote_info;
    write(res->sock, &local_info, sizeof(local_info));
    read(res->sock, &remote_info, sizeof(remote_info));
    if (modify_qp_to_init(res->qp) ||
        modify_qp_to_rtr(res->qp, remote_info.qp_num, remote_info.lid, remote_info.psn) ||
        modify_qp_to_rts(res->qp, local_info.psn)) {
        fprintf(stderr, "Failed to modify QP state\n");
        goto cleanup;
    }
    int in_fd = open(file_path, O_RDONLY);
    if (in_fd < 0) {
        fprintf(stderr, "Failed to open input file: %s\n", strerror(errno));
        goto cleanup;
    }
    ssize_t bytes_read;
    while ((bytes_read = read(in_fd, res->buf, MSG_SIZE)) > 0) {
        struct ibv_send_wr wr = {
            .wr_id = 2,
            .sg_list = &(struct ibv_sge){
                .addr = (uint64_t)res->buf,
                .length = bytes_read,
                .lkey = res->mr->lkey,
            },
            .num_sge = 1,
            .opcode = IBV_WR_SEND,
            .send_flags = IBV_SEND_SIGNALED,
        };
        struct ibv_send_wr *bad_wr;
        if (ibv_post_send(res->qp, &wr, &bad_wr)) {
            fprintf(stderr, "Failed to post send request\n");
            goto cleanup;
        }
        struct ibv_wc wc;
        while (ibv_poll_cq(res->cq, 1, &wc) == 0);
        if (wc.status != IBV_WC_SUCCESS) {
            fprintf(stderr, "Send failed with status %d\n", wc.status);
            goto cleanup;
        }
    }
    close(in_fd);
cleanup:
    if (res->qp) ibv_destroy_qp(res->qp);
    if (res->mr) ibv_dereg_mr(res->mr);
    if (res->buf) free(res->buf);
    if (res->cq) ibv_destroy_cq(res->cq);
    if (res->pd) ibv_dealloc_pd(res->pd);
    if (res->context) ibv_close_device(res->context);
    close(res->sock);
    free(res);
    return 0;
}
int main(int argc, char **argv) {
    if (argc > 2) {
        // 客户端模式
        return run_client(argv[1], argv[2]);
    } else {
        // 服务器模式
        return run_server();
    }
}