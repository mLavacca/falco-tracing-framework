#include "tracer_common.h"

int TracerCommon::open_file(int analysis_type)
{
    if(analysis_type == ONLINE_ANALYSIS)
    {
	    this->file_name = "/tmp/falco_tracer_pipe";

        remove(this->file_name.c_str());

        if(mkfifo(this->file_name.c_str(), 0666) == -1)
        {
            return -1;
        }

        this->file_pointer = open(this->file_name.c_str(), O_RDWR);
    }

    if(analysis_type == OFFLINE_ANALYSIS)
    {
        this->file_name = "/tmp/falco_tracer_file";

        remove(this->file_name.c_str());

        this->file_pointer = open(this->file_name.c_str(), O_WRONLY | O_CREAT, 0666);

        this->open_rules_file();
    }
   
    if(this->file_pointer == -1 || 
        (analysis_type == OFFLINE_ANALYSIS && this->rules_file_pointer == -1))
    {
        return -1;
    }

    return 0;
}

void TracerCommon::open_rules_file()
{
    this->file_name = "/tmp/falco_rules_names";

    remove(this->file_name.c_str());

    this->rules_file_pointer = open(this->file_name.c_str(), O_WRONLY | O_CREAT, 0666);
}

int TracerCommon::send_message(const char* message)
{
    return write(this->file_pointer, message, strlen(message));
}

int TracerCommon::write_rule_name(const char* message)
{
    return write(this->rules_file_pointer, message, strlen(message));
}

TracerCommon::~TracerCommon()
{
    close(this->file_pointer);
}
