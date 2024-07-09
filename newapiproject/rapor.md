```bash
# nodu devre dışı bırakmak
iptables -A INPUT -s 192.168.1.101 -j DROP
iptables -A OUTPUT -d 192.168.1.101 -j DROP

# nodu tekrar çalıştırmak
iptables -D INPUT -s 192.168.1.101 -j DROP
iptables -D OUTPUT -d 192.168.1.101 -j DROP

etcdctl --endpoints=http://192.168.0.1:2379,http://192.168.0.2:2379 endpoint health


# Cluster'daki nodelri incelemek ve health checking
for endpoint in $(etcdctl member list | awk '{print $4}' | cut -d',' -f1); do 
    etcdctl endpoint health --endpoints=$endpoint
done

*$ etcdutl snapshot restore snapshot.db --bump-revision 1000000000 
--mark-compacted --data-dir output-dir*

docker run -d --name etcd6 --network etcd-net1 bitnami/etcd:latest --name etcd7 --initial-advertise-peer-urls http://etcd1:2380 --listen-peer-urls http://etcd1:2380 --advertise-client-urls http://etcd1:2379 --listen-client-urls http://0.0.0.0:2379
docker run -d --name etcd7 --network etcd-net2 bitnami/etcd:latest --name etcd7 --initial-advertise-peer-urls http://etcd1:2380 --listen-peer-urls http://etcd1:2380 --advertise-client-urls http://etcd1:2379 --listen-client-urls http://0.0.0.0:2379
docker run -d --name etcd8 --network etcd-net3 bitnami/etcd:latest --name etcd8 --initial-advertise-peer-urls http://etcd2:2380 --listen-peer-urls http://etcd2:2380 --advertise-client-urls http://etcd2:2379 --listen-client-urls http://0.0.0.0:2379


docker network create etcd-net7
docker network create etcd-net8


docker run -d \
  --name etcd5\
  --net=host \
  -v /path/to/etcd/data:/etcd-data \
  quay.io/coreos/etcd:v3.5.14 \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data \
  --name etcd4 \
  --initial-advertise-peer-urls http://localhost:2385\
  --listen-peer-urls http://localhost:2385 \
  --advertise-client-urls http://localhost:2375 \
  --listen-client-urls http://localhost:2375 \
  --force-new-cluster

