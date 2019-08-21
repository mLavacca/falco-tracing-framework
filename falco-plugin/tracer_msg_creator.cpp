#include "tracer_msg_creator.h"
#include <string.h>

#define BUFFER_SIZE 100


char* TracerMsgCreator::create_message(const char* msg, const char* name, uint64_t cnt)
{
    char *buffer = new char[BUFFER_SIZE];
    char *sep = (char*)"-";
    const char *_cnt;

    memset(buffer, '\0', BUFFER_SIZE);
    strncpy(buffer, msg, strlen(msg));

    strncat(buffer, name, strlen(name));

    strncat(buffer, sep, strlen(sep));
    _cnt = std::to_string(cnt).c_str();
    strncat(buffer, _cnt, strlen(_cnt));

    buffer[strlen(buffer)] = '\n';

    return buffer;
}

char* TracerMsgCreator::create_message(const char* msg, const char* name, uint32_t tag, uint64_t cnt, uint64_t lat)
{
    char *buffer = new char[BUFFER_SIZE];
    char *sep = (char*)"-";
    const char *_tag, *_cnt, *_lat;

    memset(buffer, '\0', BUFFER_SIZE);
    strncpy(buffer, msg, strlen(msg));
    
    strncat(buffer, name, strlen(name));

    strncat(buffer, sep, strlen(sep));

    _tag = std::to_string(tag).c_str();
    strncat(buffer, _tag, strlen(_tag));

    strncat(buffer, sep, strlen(sep));

    _cnt = std::to_string(cnt).c_str();
    strncat(buffer, _cnt, strlen(_cnt));

    strncat(buffer, sep, strlen(sep));
    
    _lat = std::to_string(lat).c_str();
    strncat(buffer, _lat, strlen(_lat));

    buffer[strlen(buffer)] = '\n';

    return buffer;
}


char* TracerMsgCreator::create_message(const char* msg, uint64_t id, const char* name)
{
    char *buffer = new char[BUFFER_SIZE];
    char *sep = (char*)"-";
    const char *_id;

    memset(buffer, '\0', BUFFER_SIZE);
    strncpy(buffer, msg, strlen(msg));

    strncat(buffer, name, strlen(name));

    strncat(buffer, sep, strlen(sep));

    _id = std::to_string(id).c_str();
    strncat(buffer, _id, strlen(_id));

    buffer[strlen(buffer)] = '\n';

    return buffer;
}



char* TracerMsgCreator::create_message(const char* msg, struct timespec start_time, struct timespec end_time)
{
    char *buffer = new char[BUFFER_SIZE];
    char *sep = (char*)"-";
    const char *_start_sec, *_end_sec;

    memset(buffer, '\0', BUFFER_SIZE);
    strncpy(buffer, msg, strlen(msg));

    strncat(buffer, sep, strlen(sep));

    _start_sec = std::to_string(start_time.tv_sec).c_str();
    strncat(buffer, _start_sec, strlen(_start_sec));

    strncat(buffer, sep, strlen(sep));
    
    _end_sec = std::to_string(end_time.tv_sec).c_str();
    strncat(buffer, _end_sec, strlen(_end_sec));

    buffer[strlen(buffer)] = '\n';

    return buffer;
}

char* TracerMsgCreator::create_message(const char* msg, const char* name, uint64_t counter, uint64_t latency, const char* caller)
{
    char *buffer = new char[BUFFER_SIZE];
    char *sep = (char*)"-";
    const char *_cnt, *_lat;

    memset(buffer, '\0', BUFFER_SIZE);
    strncpy(buffer, msg, strlen(msg));
    
    strncat(buffer, name, strlen(name));

    strncat(buffer, sep, strlen(sep));

    _cnt = std::to_string(counter).c_str();
    strncat(buffer, _cnt, strlen(_cnt));

    strncat(buffer, sep, strlen(sep));
    
    _lat = std::to_string(latency).c_str();
    strncat(buffer, _lat, strlen(_lat));

    strncat(buffer, sep, strlen(sep));
    strncat(buffer, caller, strlen(caller));

    buffer[strlen(buffer)] = '\n';

    return buffer;
}