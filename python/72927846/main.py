from urllib import request
import requests
from bs4 import BeautifulSoup
import json

headers = {
    'Accept-Language': 'en-GB,en-US;q=0.9,en;q=0.8,pt;q=0.7',
    'Connection': 'keep-alive',
    'Origin': 'https://www.nationalhardwareshow.com',
    'Referer': 'https://www.nationalhardwareshow.com/',
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'cross-site',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36',
    'accept': 'application/json',
    'content-type': 'application/x-www-form-urlencoded',
    'sec-ch-ua': '".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
}
base_url='https://api.reedexpo.com/v1/organisations/'

params = {
    'x-algolia-agent': 'Algolia for vanilla JavaScript 3.27.1',
    'x-algolia-application-id': 'XD0U5M6Y4R',
    'x-algolia-api-key': 'd5cd7d4ec26134ff4a34d736a7f9ad47',
}

data = '{"params":"query=&page=0&facetFilters=&optionalFilters=%5B%5D"}'

resp = requests.post('https://xd0u5m6y4r-3.algolianet.com/1/indexes/event-edition-eve-e6b1ae25-5b9f-457b-83b3-335667332366_en-us/query', params=params, headers=headers, data=data).json()
productlinks=[]
for item in resp['hits']:
    url=base_url+item['organisationGuid']+"/exhibiting-organisations?eventEditionId="+item['eventEditionExternalId']
    productlinks.append(url)
    
for link in productlinks:
    print("link", link)
    title=link['_embedded']['companyName'].json()
    print(title)