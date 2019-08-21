#include <stdlib.h>
#include <assert.h>

#include "tracer.h"
#include "tracer_interface.h"


Tracer *global_tracer;
uint64_t rule_begin_time;


/*
 * signal handlers
 */
void send_data(int signum)
{
	global_tracer->send_summary();
}


void flush_data(int signum)
{
    global_tracer->flush_stats();
}


/*
 * Tracer implementation
 */
Tracer::Tracer()
{
    printf("TRACER - tracer activity started\n");

    this->last_counter = -1;

    this->timestamp_array_func = new(std::nothrow) uint64_t[N_MONITORED_FUNCTIONS];
    this->timestamp_array_rule = new(std::nothrow) uint64_t[N_RULES];

    for(int i = 0; i < N_RULES; i++)
    {
        this->rule_array[i].counter = 0;
        this->rule_array[i].total_time = 0;
    }

    for(int i = 0; i < N_MONITORED_COUNTERS; i++)
    {
        for(int j = 0; j < N_MONITORED_FUNCTIONS; j++)
        {
            this->func_array[i][j].counter = 0;
            this->func_array[i][j].total_time = 0;
        }
    }

    for(int i = 0; i < N_MONITORED_FUNCTIONS; i++)
    {
        this->func_last_call_array[i].counter = 0;
        this->func_last_call_array[i].total_time = 0;
    }
}


int Tracer::configure_tracer(int analysis_type)
{
    this->analysis_type = analysis_type;

    if(analysis_type == ONLINE_ANALYSIS)
    {
        if(this->configure_signals_callbacks() == -1)
        {
            return -1;
        }
    }
    
    if(this->open_file(analysis_type) == -1)
    {
        return -1;
    }

    return 0;
}

int Tracer::configure_signals_callbacks()
{
	if (signal(SEND_STATS, send_data) == SIG_ERR ||
		signal(FLUSH_DATA, flush_data) == SIG_ERR)
	{
		return -1;
	}

    return 0;
}


int Tracer::set_func_start(int index)
{
    if(index < 0 || index >= N_MONITORED_FUNCTIONS)
    {
        return -1;
    }

    this->timestamp_array_func[index] = _rdtsc_begin();

    return 0;
}

int Tracer::set_func_end(int index, int caller)
{
    uint64_t end = _rdtsc_end();
    uint64_t latency = end - this->timestamp_array_func[index];

    if(index < 0 || index >= N_MONITORED_FUNCTIONS)
    {
        return -1;
    }

    this->func_last_call_array[index].counter ++;
    this->func_last_call_array[index].total_time += latency;

    if(caller == ROOT)
    {
        this->func_last_call_array[index].caller = "root";
    }
    else
    {
        this->func_last_call_array[index].caller = function_names[caller];
    }
    
    return 0;
}


int Tracer::set_func_end_cont(int index, int triggering_counter, int caller)
{
    uint64_t end = _rdtsc_end();
    uint64_t latency = end - this->timestamp_array_func[index];

    if(index < 0 || index >= N_MONITORED_FUNCTIONS)
    {
        return -1;
    }

    if(triggering_counter < 0 || triggering_counter >= N_MONITORED_COUNTERS)
    {
        return -1;
    }

    if(caller < -1 || caller >= N_MONITORED_FUNCTIONS)
    {
        return -1;
    }

    if(this->last_counter == -1)
    {
        this->last_counter = triggering_counter;
    }

    this->func_last_call_array[index].counter ++;
    this->func_last_call_array[index].total_time += latency;

    if(caller == ROOT)
    {
        this->func_last_call_array[index].caller = "root";
    }
    else
    {
        this->func_last_call_array[index].caller = function_names[caller];
    }
    
    if(triggering_counter == CYCLE_COMPLETED_CNT  || triggering_counter == SCAP_TIMEOUT_CNT)
    {
        for(int j = 0; j < N_MONITORED_FUNCTIONS; j++)
        {
            this->func_array[this->last_counter][j].counter += this->func_last_call_array[j].counter;
            this->func_array[this->last_counter][j].total_time += this->func_last_call_array[j].total_time;
            this->func_array[this->last_counter][j].caller = this->func_last_call_array[j].caller;

            this->func_last_call_array[j] = {.counter = 0, .total_time = 0, .caller = ""};
        }

        this->last_counter = -1;
    }

    return 0;
}


int Tracer::set_rule_start()
{
    rule_begin_time = _rdtsc_begin();

    return 0;
}


int Tracer::set_rule_end(int check_id, int tag_id, bool broken)
{
    uint64_t end = _rdtsc_end();

    assert(end > rule_begin_time);

    if(check_id < 0 || check_id >= N_RULES)
    {
        return -1;
    }

    if(broken)
    {
        this->rule_array_broken[check_id].counter ++;
        this->rule_array_broken[check_id].total_time += (end - rule_begin_time);
        this->rule_array_broken[check_id].tag = tag_id;
    }
    else
    {
        this->rule_array[check_id].counter ++;
        this->rule_array[check_id].total_time += (end - rule_begin_time);
        this->rule_array[check_id].tag = tag_id;
    }  

    return 0;
}


void Tracer::start_rules_names_sending()
{
    clock_gettime(CLOCK_REALTIME, &start_time);
    
    if(this->analysis_type == ONLINE_ANALYSIS)
    {
        this->send_message("TRACER INFO-START RULES NAMES\n");
    }

    if(this->analysis_type == OFFLINE_ANALYSIS)
    {
        this->write_rule_name("TRACER INFO-START RULES NAMES\n");
    }
}


void Tracer::end_rules_names_sending()
{
    if(this->analysis_type == ONLINE_ANALYSIS)
    {
        this->send_message("TRACER INFO-END RULES NAMES\n");
    }

    if(this->analysis_type == OFFLINE_ANALYSIS)
    {
        this->write_rule_name("TRACER INFO-END RULES NAMES\n");
    }
}


int Tracer::send_rule_name(int check_id, const char* name)
{
    if(check_id < 0 || check_id >= N_RULES)
    {
        return -1;
    }

    char* message = this->create_message(
            "TRACER-", check_id, name);

    if(this->analysis_type == ONLINE_ANALYSIS)
    {
        this->send_message((const char*)message);
    }

    if(this->analysis_type == OFFLINE_ANALYSIS)
    {
        this->write_rule_name((const char*)message);
    }
    

    return 0;
}


int Tracer::inc_counter(int index)
{
    if(index < 0 || index >= N_MONITORED_FUNCTIONS)
    {
        return -1;
    }

    this->counter_array[index].counter ++;

    return 0;
}


void Tracer::send_summary()
{   
    clock_gettime(CLOCK_REALTIME, &end_time);

    char* message = this->create_message("TRACER INFO-START SUMMARY", start_time, end_time);
    this->send_message(message);

    this->send_message("TRACER INFO-START STACKTRACES\n");
    uint32_t tag;
    uint64_t cnt, lat;
    const char* caller;
    for(int i = 0; i < N_MONITORED_COUNTERS; i++)
    {
        if(this->func_array[i][0].counter == 0)
        {
            continue;
        }
        char* message = this->create_message(
                "TRACER INFO-START STACKTRACE-", (const char*)this->counter_names[i], this->counter_array[i].counter);

        this->send_message((const char*)message); 

        for(int j = 0; j < N_MONITORED_FUNCTIONS; j++)
        {
            lat = 0;
            cnt = this->func_array[i][j].counter;
            caller = this->func_array[i][j].caller;

            if(cnt == 0)
            {
                continue;
            }

            lat = this->func_array[i][j].total_time / cnt;
        
            char* message = this->create_message(
                "TRACER-", (const char*)this->function_names[j], cnt, lat, caller);

            this->send_message((const char*)message);
        }

        this->send_message("TRACER INFO-END STACKTRACE\n"); 
    }
    this->send_message("TRACER INFO-END STACKTRACES\n");

    this->send_message("TRACER INFO-START COUNTERS\n");
    for(int i = 0; i < N_MONITORED_COUNTERS; i++)
    {
        cnt = this->counter_array[i].counter;
        if(cnt  == 0)
        {
            continue;
        }
        
        char* message = this->create_message(
            "TRACER-", (const char*)this->counter_names[i], cnt);

        this->send_message((const char*)message);   
    }
    this->send_message("TRACER INFO-END COUNTERS\n");

    this->send_message("TRACER INFO-START UNBROKEN RULES\n");
    for(int i = 0; i < N_RULES; i++)
    {
        cnt = this->rule_array[i].counter;
        tag = this->rule_array[i].tag;
        uint64_t total_time = this->rule_array[i].total_time;

        if(cnt != 0)
        {
            lat = total_time / cnt;
            char* message = this->create_message(
                "TRACER-", std::to_string(i).c_str(), tag, cnt, lat);

            this->send_message((const char*)message);
        }
    }
    this->send_message("TRACER INFO-END UNBROKEN RULES\n");

    this->send_message("TRACER INFO-START BROKEN RULES\n");
    for(int i = 0; i < N_RULES; i++)
    {
        cnt = this->rule_array_broken[i].counter;
        tag = this->rule_array_broken[i].tag;
        uint64_t total_time = this->rule_array_broken[i].total_time;

        if(cnt != 0)
        {
            lat = total_time / cnt;
            char* message = this->create_message(
                "TRACER-", std::to_string(i).c_str(), tag, cnt, lat);

            this->send_message((const char*)message);
        }
    }
    this->send_message("TRACER INFO-END BROKEN RULES\n");
    
    this->send_message("TRACER INFO-END SUMMARY\n");
}


void Tracer::flush_stats()
{
    for(int i = 0; i< N_MONITORED_COUNTERS; i++)
    {
        this->counter_array[i].counter = 0;
    }

    for(int i = 0; i < N_MONITORED_COUNTERS; i++)
    {
        for(int j = 0; j < N_MONITORED_FUNCTIONS; j++)
        {
            this->func_array[i][j].counter = 0;
            this->func_array[i][j].total_time = 0;
        } 
    }

    for(int i = 0; i< N_RULES; i++)
    {
        this->rule_array[i].counter = 0;
        this->rule_array[i].total_time = 0;

        this->rule_array_broken[i].counter = 0;
        this->rule_array_broken[i].total_time = 0;
    }

    clock_gettime(CLOCK_REALTIME, &start_time);
}


Tracer::~Tracer()
{    
    printf("TRACER - tracer activity ended\n");
}