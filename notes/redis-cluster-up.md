

———— OK ————
docker run -d --name redis-node1 -p 7000:7000 redis:latest redis-server --port 7000 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000
docker run -d --name redis-node2 -p 7001:7001 redis:latest redis-server --port 7001 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000
docker run -d --name redis-node3 -p 7002:7002 redis:latest redis-server --port 7002 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000

docker run -d --name redis-node4 -p 7003:7003 redis:latest redis-server --port 7003 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000
docker run -d --name redis-node5 -p 7004:7004 redis:latest redis-server --port 7004 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000
docker run -d --name redis-node6 -p 7005:7005 redis:latest redis-server --port 7005 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000
———— OK ————


# Check IPs of the nodes and Run the command to create the Redis Cluster with the new nodes


docker inspect redis-node1 | grep IPAddress
            "SecondaryIPAddresses": null,
            "IPAddress": "192.168.215.2",
                    "IPAddress": "192.168.215.2",

docker inspect redis-node2 | grep IPAddress
            "SecondaryIPAddresses": null,
            "IPAddress": "192.168.215.3",
                    "IPAddress": "192.168.215.3",

.....


redis-cli --cluster create \
192.168.215.2:7000 192.168.215.3:7001 192.168.215.4:7002 \
192.168.215.5:7003 192.168.215.6:7004 192.168.215.7:7005 \
--cluster-replicas 1



Finally conf @ main:

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.215.2:7000",
			"192.168.215.3:7001",
			"192.168.215.4:7002",
			"192.168.215.5:7003",
			"192.168.215.6:7004",
			"192.168.215.7:7005",
		},
	})