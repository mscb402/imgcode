import requests
import re
##import sys




def uploadFile(filename):
    url='http://image.baidu.com/pictureup/uploadshitu'
    files={'image':'' }
    imgage = open(filename,'rb')
    files['image']=imgage
    c=requests.post(url,files=files)
    geturl1=re.sub('%3A',':',re.findall('queryImageUrl=(.*?)&querySign',c.url)[0])
    realurl=re.sub('%2F','/',geturl1)
    return realurl


print(uploadFile(r'test.png'))

