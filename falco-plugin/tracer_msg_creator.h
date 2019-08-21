#pragma once

#include <string>


class TracerMsgCreator
{
public:
    TracerMsgCreator() {};
    ~TracerMsgCreator() {};
protected:
    char* create_message(const char* msg, const char* name, uint64_t cnt);
    char* create_message(const char* msg, uint64_t id, const char* name);
    char* create_message(const char* msg, struct timespec start_time, struct timespec end_time);
    char* create_message(const char* msg, const char* name, uint32_t tag, uint64_t cnt, uint64_t lat);
    char* create_message(const char* msg, const char* name, uint64_t counter, uint64_t latency, const char* caller);
};