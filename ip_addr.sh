ip_addr=$(hostname -I)
echo $ip_addr
s3cmd del -r s3://$1/ip_addr
touch $ip_addr
s3cmd put ./$ip_addr s3://$1/ip_addr/$ip_addr