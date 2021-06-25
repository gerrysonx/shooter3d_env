import boto3
from boto3.session import Session

mybucket = "haytham_traindata__datamingplatform"

access_key = "vkimgpull-datamingplatform-f98d2d53"

secret_key = "vkimgpull-datamingplatform-6ab40a41"

url = "http://shpublicrgw.cephrados.so.db:7480"

session = Session(access_key, secret_key)

s3_client = session.client('s3', endpoint_url=url)

s3 = session.resource('s3', endpoint_url=url)

bucket = s3.Bucket('haytham_traindata__datamingplatform')

objectlist =  bucket.objects.filter(Prefix = 'distribute_collected_train_data')
for obj in objectlist:
    s3.Object(mybucket, obj.key).delete()