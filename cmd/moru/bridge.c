typedef void (*log_write_t)(const void *msg, int len);

void bridge_log_write(log_write_t f, const void *msg, int len) {
    f(msg, len);
}
