#pragma once

/*
 * type of analysis
 */
#define ONLINE_ANALYSIS 1
#define OFFLINE_ANALYSIS 2

/*
 * signals
 */
#define SEND_STATS 34
#define FLUSH_DATA 35
#define SEND_RULES_NAMES 36

/*
 * monitored functions
 */
#define N_MONITORED_FUNCTIONS 9

#define ANALYSIS_CYCLE 0
#define SINSP_NEXT 1
#define STATS_FILE_WRITER_HANDLE 2
#define SYSCALL_EVT_DROP_MGR_PROCESS_EVENT 3
#define EV_FALCO_CONSIDER 4
#define PROCESS_SINSP_EVENT 5
#define FALCO_ENGINE_SHOULD_DROP_EVT 6
#define RULESET_FILTERS_RUN 7
#define FALCO_OUTPUTS_HANDLE_EVENT 8

#define ROOT -1

/*
 * monitored counters
 */
#define N_MONITORED_COUNTERS 5

#define SCAP_TIMEOUT_CNT 0
#define EVENT_NOT_CONSIDERED_CNT 1
#define EVENT_DROPPED_CNT 2
#define RULES_NOT_BROKEN_CNT 3
#define CYCLE_COMPLETED_CNT 4

/*
 * maximum number of rules
 */
#define N_RULES 1024
