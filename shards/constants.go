package shards

// These are the endpoints required for
// cheshire sharding to work.
const (
	//router table get endpoint
	// response format
	// {
	//  "strest" :{...}
	//  "router_table" : <the router table>
	// }
	// @method GET
	ROUTERTABLE_GET = "/__c/rt/get"

	// Sets the router table on this server
	// @method POST
	// @param router_table The router table
	ROUTERTABLE_SET = "/__c/rt/set"

	// @method POST
	// @param partition
	PARTITION_LOCK   = "/__c/pt/lock"
	
	// @method POST
	// @param partition
	PARTITION_UNLOCK = "/__c/pt/unlock"

	// Delete a partition from this server
	// @method DELETE
	// @param partition
	PARTITION_DELETE = "/__c/pt/delete"

	// Is a ping endpoint to check for liveness and
	// to check the revision of the router table.
	// response format
	// {
	//  "strest" :{...}
	//  "ts" : <ISOFORMATED TIMESTAMP>
	//  "rt_revision" : <router table revision>
	// }
	// @method GET
	//
	CHECKIN = "/__c/checkin"

	// Creates a stream of data for the given partition
	// @param partition the int partition
	// @method GET
	PARTITION_EXPORT = "/__c/pt/export"

	// Initializes an import request between two shards
	//
	// @method POST
	// @param partition the partition to import data
	// @param source the http address to pull data from in the form http://address:port
	PARTITION_IMPORT = "/__c/pt/import"
)

//These are the required return error codes for various situations
const (
	// We reserve 630-640 for router table issues

	// return when the requester has an old router table
	E_ROUTER_TABLE_OLD = 632

	// requester has a newer router table then us, request they update ours
	E_SEND_ROUTER_TABLE = 633

	// the requested partition is locked.  requester should try back in a bit
	E_PARTITION_LOCKED = 634

	// The requested partition does not live on this shard
	E_NOT_MY_PARTITION = 635
)

// Param Names
const (
	//The partition val (an integer from 0 to TotalPartitions)
	P_PARTITION = "_p"

	// The version of the router table
	P_REVISION = "_v"

	//The shard key, should only be used when passing to a proxy
	P_SHARD_KEY = "_sk"
	
	//The query type.
	// This defines how the request can be handled by the router.
	// Possible values:
	// single : return a single result (the first response received)
	// all : (default) return values for all servers, will make an effort to retry on failure, but will generally return error results.
	// all_q : return values for all servers (queue requests if needed, retry until response).  This would typically be for posting
	// none_q : returns success immediately, queues the request and make best effort to ensure it is delivered (TODO)
	P_QUERY_TYPE = "_qt"
)
