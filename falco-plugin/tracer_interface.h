#pragma once
#include "defined_values.h"

#ifdef __cplusplus
extern "C"{
#endif

typedef struct rule_info{
    char* name;
    char* type;
}rule_info;

int init_tracer();

void set_analysis_type(int type);
int get_analysis_type();

void set_current_tag(int tag);
void set_current_rule(int id);
int set_rule_start();
int set_rule_end(bool broken);

void start_rules_names_sending();
int send_rule_name(int check_id, const char* name);
void end_rules_names_sending();

int set_func_start(int i);
int set_func_end(int i, int caller);
int set_func_end_cont(int i, int counter, int caller);

void write_scap_file_stats();

void reset_last_counter();

int inc_counter(int i);

void close_tracer();

#ifdef __cplusplus
}
#endif

