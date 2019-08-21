#include "tracer_interface.h"
#include "tracer.h"
#include "tracer_common.h"
#include <fstream>

extern Tracer *global_tracer;
uint32_t current_rule;
uint32_t current_tag;
int analysis_type = ONLINE_ANALYSIS;


int init_tracer()
{
    global_tracer = new(std::nothrow) Tracer();
    if(global_tracer == NULL)
    {
        return -1;
    }

    return global_tracer->configure_tracer(analysis_type);
}

void set_analysis_type(int type)
{
    analysis_type = type;
}

int get_analysis_type()
{
    return analysis_type;
}

void set_current_tag(int tag)
{
    current_tag = tag;
}

void set_current_rule(int id)
{
    current_rule = id;
}

int set_rule_start()
{
    return global_tracer->set_rule_start();
}

int set_rule_end(bool broken)
{
    return global_tracer->set_rule_end(current_rule, current_tag, broken);
}

void start_rules_names_sending()
{
    global_tracer->start_rules_names_sending();
}

void end_rules_names_sending()
{
    global_tracer->end_rules_names_sending();
}

int send_rule_name(int check_id, const char* name)
{
    return global_tracer->send_rule_name(check_id, name);
}

int set_func_start(int i)
{
    return global_tracer->set_func_start(i);
}

int set_func_end(int i, int caller)
{
    return global_tracer->set_func_end(i, caller);
}

int set_func_end_cont(int i, int counter, int caller)
{
    return global_tracer->set_func_end_cont(i, counter, caller);
}

void write_scap_file_stats()
{
    global_tracer->send_summary();
}

int inc_counter(int i)
{
    return global_tracer->inc_counter(i);
}

void reset_last_counter()
{
    global_tracer->last_counter = -1;
}

void close_tracer()
{
    delete global_tracer;
}