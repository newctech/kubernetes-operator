#!/bin/bash

[ -e /etc/init.d/functions ] && . /etc/init.d/functions || exit
[ -e ./config.sh ] && . ./config.sh || exit

KUBE_MASTER_BIN_DIR="../bin/${KUBERNETES_VER}"
KUBE_MASTER_CONFIG_DIR="../config/master"
KUBE_MASTER_SYSTEMD_CONFIG_DIR="../systemd"
GENERATE_CERTS_FILE="../certs/master/gen_cert.sh"
GENERATE_KUBECONFIG_FILE="../kubeconfig/generate_master_kubeconfig.sh"

DEST_CONFIG_DIR="/etc/kubernetes"
DEST_SYSTEMD_DIR="/usr/lib/systemd/system"
DEST_CERTS_DIR="/etc/kubernetes/ssl"
KUBE_MASTER_LOG="/var/log/kubernetes"

cp ${KUBE_MASTER_BIN_DIR}/{kube-apiserver,kube-controller-manager,kube-scheduler,kubectl} /usr/bin/
cp ${KUBE_MASTER_SYSTEMD_CONFIG_DIR}/{kube-apiserver.service,kube-controller-manager.service,kube-scheduler.service} ${DEST_SYSTEMD_DIR}/

# cp config, apiserver config controller-manager scheduler 
cp ${KUBE_MASTER_CONFIG_DIR}/* ${DEST_CONFIG_DIR}/

# config etcd server
etcd_num=$(echo ${ETCD_HOSTS} | awk -F ',' '{print NF}')
etcd_cluster=""
for i in `seq 1 ${etcd_num}`;do
	ip=$(echo ${ETCD_HOSTS} | awk -v idx=$i -F ',' '{print $idx}')
    cluster=$(echo "https://${ip}:2379")
    if [ $i -ne ${etcd_num} ];then
        cluster="${cluster},"
    fi
    etcd_cluster="${etcd_cluster}${cluster}"
done

# update config
sed -i -e "s#--etcd-servers=<etcd_cluster>#--etcd-servers=${etcd_cluster}#g" ${DEST_CONFIG_DIR}/apiserver
sed -i -e "s#--master=https://<apiserver_ip>:6443#--master=https://${LOCAL_IP}:6443#g" ${DEST_CONFIG_DIR}/config

# generate ssl 
bash ${GENERATE_CERTS_FILE}
[ $? -eq 0 ] && echo "generate certs success" || exit 1
cp ${GENERATE_CERTS_FILE}/output/* ${DEST_CERTS_DIR}/

# generate kubeconfig
sed -i -e "s#https://<apiserver_ip>:6443#https://${LOCAL_IP}:6443#g" ${GENERATE_KUBECONFIG_FILE}
bash ${GENERATE_KUBECONFIG_FILE}
[ $? -eq 0 ] && echo "generate kubeconfig success" || exit 1
cp ${GENERATE_KUBECONFIG_FILE}/output/* ${DEST_CONFIG_DIR}/

# mkdir master log dir 
[ -d ${KUBE_MASTER_LOG} ] || mkdir -pv ${KUBE_MASTER_LOG}

# start service
systemctl daemon-reload
systemctl enable kube-apiserver kube-controller-manager kube-scheduler 
systemctl start kube-apiserver kube-controller-manager kube-scheduler
systemctl status kube-apiserver kube-controller-manager kube-scheduler
