#pragma once 

#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <csignal>
#include <string>
#include <iostream>

#include "tsc.h"
#include "defined_values.h"
#include "tracer_common.h"
#include "tracer_msg_creator.h"

class Tracer : public TracerCommon, public TracerMsgCreator{

private:

struct timespec start_time;
struct timespec end_time;

typedef struct func_stats{
    uint64_t counter;
    uint64_t total_time;
    const char* caller;
}func_stats;

typedef struct cnt_stats{
    uint64_t counter;
}cnt_stats;

typedef struct rule_stats{
    uint32_t tag;
    uint64_t counter;
    uint64_t total_time;
}rule_stats;

    const char* function_names[N_MONITORED_FUNCTIONS] = {
        "analysis_cycle",
        "Sinsp::next",
        "StatsFileWriter::handle",
        "syscall_evt_drop_mgr::process_event",
        "sinsp_evt::falco_consider",
        "falco_engine::process_sinsp_event",
        "falco_engine::should_drop_evt",
        "falco_sinsp_ruleset::run"
    };

    const char* counter_names[N_MONITORED_COUNTERS] = {
        "scap_timeout",
        "event_not_considered",
        "event_dropped",
        "rules_unbroken",
        "rules_broken",
    };

    func_stats func_last_call_array[N_MONITORED_FUNCTIONS];
    func_stats func_array[N_MONITORED_COUNTERS][N_MONITORED_FUNCTIONS];
    cnt_stats counter_array[N_MONITORED_COUNTERS];

    rule_stats rule_array[N_RULES];
    rule_stats rule_array_broken[N_RULES];
    
    uint64_t *timestamp_array_func;
    uint64_t *timestamp_array_rule;

    int analysis_type;

public:
    int last_counter;

    int configure_tracer(int analysis_type);
    
    Tracer();
    int configure_signals_callbacks();
    
    int set_func_start(int index);
    int set_func_end(int index, int caller);
    int set_func_end_cont(int index, int triggering_counter, int caller);

    int set_rule_start();
    int set_rule_end(int check_id, int tag_id, bool broken);

    void start_rules_names_sending();
    int send_rule_name(int check_id, const char* name);
    void end_rules_names_sending();

    int inc_counter(int index);

    void send_summary();

    void write_stats_on_file();

    void flush_stats();
    
    ~Tracer();
};