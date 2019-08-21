#include "tsc.h"

uint64_t a = 0;
uint64_t d = 0;

__uint64_t _rdtsc_begin(){
    // a: low bits; d: high bits
    a = 0;
    d = 0;

    asm volatile("cpuid" : : : "%rax", "%rbx", "%rcx", "%rdx", "memory");
    asm volatile("rdtsc" : "=a" (a), "=d" (d) : : "%rbx", "%rcx", "memory");

    return a | (d << 32);
}

__uint64_t _rdtsc_end(){
    // a: low bits; d: high bits
    a = 0;
    d = 0;

    asm volatile("rdtscp" : "=a" (a), "=d" (d) : : "%rbx", "%rcx", "memory");
    asm volatile("cpuid" : : : "%rax", "%rbx", "%rcx", "%rdx", "memory");

    return a | (d << 32);
}
