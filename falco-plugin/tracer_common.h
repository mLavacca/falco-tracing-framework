#pragma once

#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <csignal>
#include <string>
#include <exception>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h> 

#include "defined_values.h"
#include "tracer_interface.h"

class TracerCommon
{
private:
    int file_pointer;
    int rules_file_pointer;
    std::string file_name;
    std::string rules_file_name;

    void open_rules_file();
public:
    TracerCommon() {};
    int open_file(int analysis_type);
    int close_file();
    ~TracerCommon();
protected:
    int send_message(const char* message);
    int write_rule_name(const char* message);
};