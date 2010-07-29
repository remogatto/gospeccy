// Type and constant definitions to be processed by the 'godefs' tool.
// This file is basically a customized partial copy of <linux/perf_event.h>.

#include <linux/unistd.h>	// __NR_perf_event_open
#include <linux/perf_event.h>	// Requires Linux kernel >= 2.6.33

enum {
	$SYS_PERF_OPEN = __NR_perf_event_open
};

enum /*perf_type_id*/ {
	$TYPE_HARDWARE		= PERF_TYPE_HARDWARE,
	$TYPE_SOFTWARE		= PERF_TYPE_SOFTWARE,
	$TYPE_TRACEPOINT	= PERF_TYPE_TRACEPOINT,
	$TYPE_HW_CACHE		= PERF_TYPE_HW_CACHE,
	$TYPE_RAW		= PERF_TYPE_RAW,
	$TYPE_BREAKPOINT	= PERF_TYPE_BREAKPOINT
};

enum /*perf_hw_id*/ {
	$HW_CPU_CYCLES		= PERF_COUNT_HW_CPU_CYCLES,
	$HW_INSTRUCTIONS	= PERF_COUNT_HW_INSTRUCTIONS,
	$HW_CACHE_REFERENCES	= PERF_COUNT_HW_CACHE_REFERENCES,
	$HW_CACHE_MISSES	= PERF_COUNT_HW_CACHE_MISSES,
	$HW_BRANCH_INSTRUCTIONS	= PERF_COUNT_HW_BRANCH_INSTRUCTIONS,
	$HW_BRANCH_MISSES	= PERF_COUNT_HW_BRANCH_MISSES,
	$HW_BUS_CYCLES		= PERF_COUNT_HW_BUS_CYCLES
};

enum /*perf_hw_cache_id*/ {
	$HW_CACHE_L1D	= PERF_COUNT_HW_CACHE_L1D,
	$HW_CACHE_L1I	= PERF_COUNT_HW_CACHE_L1I,
	$HW_CACHE_LL	= PERF_COUNT_HW_CACHE_LL,
	$HW_CACHE_DTLB	= PERF_COUNT_HW_CACHE_DTLB,
	$HW_CACHE_ITLB	= PERF_COUNT_HW_CACHE_ITLB,
	$HW_CACHE_BPU	= PERF_COUNT_HW_CACHE_BPU
};

enum /*perf_hw_cache_op_id*/ {
	$HW_CACHE_OP_READ	= PERF_COUNT_HW_CACHE_OP_READ,
	$HW_CACHE_OP_WRITE	= PERF_COUNT_HW_CACHE_OP_WRITE,
	$HW_CACHE_OP_PREFETCH	= PERF_COUNT_HW_CACHE_OP_PREFETCH
};

enum /*perf_hw_cache_op_result_id*/ {
	$HW_CACHE_RESULT_ACCESS	= PERF_COUNT_HW_CACHE_RESULT_ACCESS,
	$HW_CACHE_RESULT_MISS	= PERF_COUNT_HW_CACHE_RESULT_MISS
};

enum /*perf_sw_ids*/ {
	$SW_CPU_CLOCK		= PERF_COUNT_SW_CPU_CLOCK,
	$SW_TASK_CLOCK		= PERF_COUNT_SW_TASK_CLOCK,
	$SW_PAGE_FAULTS		= PERF_COUNT_SW_PAGE_FAULTS,
	$SW_CONTEXT_SWITCHES	= PERF_COUNT_SW_CONTEXT_SWITCHES,
	$SW_CPU_MIGRATIONS	= PERF_COUNT_SW_CPU_MIGRATIONS,
	$SW_PAGE_FAULTS_MIN	= PERF_COUNT_SW_PAGE_FAULTS_MIN,
	$SW_PAGE_FAULTS_MAJ	= PERF_COUNT_SW_PAGE_FAULTS_MAJ,
	$SW_ALIGNMENT_FAULTS	= PERF_COUNT_SW_ALIGNMENT_FAULTS,
	$SW_EMULATION_FAULTS	= PERF_COUNT_SW_EMULATION_FAULTS
};

// Same as the struct 'perf_event_attr', but without the unions and bitfields.
struct GO__perf_event_attr {
	__u32			type;
	__u32			size;
	__u64			config;

	__u64			sample_periodOrFreq;

	__u64			sample_type;
	__u64			read_format;

	__u64			flags;

	__u32			wakeup_eventsOrWatermark;

	__u32			bp_type;
	__u64			bp_addr;
	__u64			bp_len;
};

typedef struct GO__perf_event_attr $Attr;

enum {
	$ATTR_SIZE = sizeof(struct GO__perf_event_attr)
};

enum {
	$FLAG_DISABLED 			= (1 << 0),
	$FLAG_INHERIT	 		= (1 << 1),
	$FLAG_PINNED			= (1 << 2),
	$FLAG_EXCLUSIVE			= (1 << 3),
	$FLAG_EXCLUDE_USER		= (1 << 4),
	$FLAG_EXCLUDE_KERNEL	= (1 << 5),
	$FLAG_EXCLUDE_HV		= (1 << 6),
	$FLAG_EXCLUDE_IDLE		= (1 << 7),
	$FLAG_MMAP				= (1 << 8),
	$FLAG_COMM				= (1 << 9),
	$FLAG_FREQ				= (1 << 10),
	$FLAG_INHERIT_STAT		= (1 << 11),
	$FLAG_ENABLE_ON_EXEC	= (1 << 12),
	$FLAG_TASK				= (1 << 13),
	$FLAG_WATERMARK			= (1 << 14)
};
