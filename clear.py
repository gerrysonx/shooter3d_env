import boto3
from boto3.session import Session
import sys
import json

class clear:
    @staticmethod
    def clear_model(idx):
        mybucket = idx

        access_key = "vkimgpull-datamingplatform-f98d2d53"

        secret_key = "vkimgpull-datamingplatform-6ab40a41"

        url = "http://shlightspeedrgw.cephrados.so.db:7480"

        session = Session(access_key, secret_key)

        s3_client = session.client('s3', endpoint_url=url)

        s3 = session.resource('s3', endpoint_url=url)

        bucket = s3.Bucket(mybucket)

        objectlist =  bucket.objects.filter(Prefix = 'model')
        for obj in objectlist:
            if not 'model/model_0001' in obj.key:
                s3.Object(mybucket, obj.key).delete()
        objectlist =  bucket.objects.filter(Prefix = 'ckpt')
        for obj in objectlist:
            if not 'ckpt/model_0001' in obj.key:
                s3.Object(mybucket, obj.key).delete()
        objectlist =  bucket.objects.filter(Prefix = 'ckpt')
        for obj in objectlist:
            copy_source = {
                'Bucket': mybucket,
                'Key': obj.key
            }
            s3.meta.client.copy(copy_source, mybucket, 'ckpt/model_in_train/{}'.format(obj.key.split('/')[-1]))
        model_info={"Model": ["0001"], "Score": [1000]}
        with open("./model_list.json","w") as f:
            json.dump(model_info,f)
        response = s3.meta.client.upload_file('./model_list.json', mybucket, 'model/model_list.json', ExtraArgs={'ACL':'public-read'})
        if response != None:
            print(response)

    @staticmethod
    def clear_data():
        mybucket = "haytham_traindata__datamingplatform"

        if len(sys.argv) > 1:
            mybucket = sys.argv[1]

        access_key = "vkimgpull-datamingplatform-f98d2d53"

        secret_key = "vkimgpull-datamingplatform-6ab40a41"

        url = "http://shlightspeedrgw.cephrados.so.db:7480"

        session = Session(access_key, secret_key)

        s3_client = session.client('s3', endpoint_url=url)

        s3 = session.resource('s3', endpoint_url=url)

        bucket = s3.Bucket(mybucket)

        objectlist =  bucket.objects.filter(Prefix = 'distribute_collected_train_data')
        for obj in objectlist:
            s3.Object(mybucket, obj.key).delete()
    