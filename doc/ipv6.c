#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main()
{
    FILE *fp;
    char line[512];
    char *src_pos;
    char ipv6_addr[128] = {0};

    // 执行命令
    fp = popen("ip -6 route get 2001:4860:4860::8888 2>/dev/null", "r");
    if (fp == NULL) {
        perror("popen failed");
        return 1;
    }

    // 读取输出并查找 src
    while (fgets(line, sizeof(line), fp) != NULL) {
        // 查找 "src " 关键字
        if ((src_pos = strstr(line, "src ")) != NULL) {
            src_pos += 4;  // 跳过 "src "

            // 提取 IPv6 地址（直到空格或换行）
            sscanf(src_pos, "%127s", ipv6_addr);

            printf("%s\n", ipv6_addr);
            break;
        }
    }

    pclose(fp);

    if (strlen(ipv6_addr) == 0) {
        fprintf(stderr, "未找到源 IPv6 地址\n");
        return 1;
    }

    return 0;
}
