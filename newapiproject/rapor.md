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


