import requests

from dataclasses import dataclass

from icecream import ic


@dataclass
class Keys:
    private: str
    public: str


class CypherClient:
    def __init__(self, base_url: str):
        self.base_url = base_url

    def create_keys(self) -> Keys:
        resp = requests.get(url=self.base_url + '/api/cypher/elliptic/keys')
        resp.raise_for_status()

        resp_json = resp.json()

        return Keys(resp_json['private'], resp_json['public'])

    def encrypt(self, text: str, public_key: str) -> bytes:
        resp = requests.post(url=self.base_url + '/api/cypher/elliptic/encrypt', files={
            'text': ('', text),
            'pemKey': ('public_key.pem', public_key),
        })
        resp.raise_for_status()

        return resp.content

    def decrypt(self, encrypted_text: bytes, private_key: str) -> str:
        files = {
            'pemKey': ('private_key.pem', private_key),
            'encryptedData': ('encrypted_data.bin', encrypted_text)
        }

        resp = requests.post(url=self.base_url + '/api/cypher/elliptic/decrypt',  files=files)
        resp.raise_for_status()

        resp_json = resp.json()

        return resp_json['text']


def main():
    base_url = "http://localhost"

    client = CypherClient(base_url)
    keys = client.create_keys()
    ic(keys)
    encrypted_text = client.encrypt("hello", public_key=keys.public)
    ic(encrypted_text)

    text = client.decrypt(encrypted_text, private_key=keys.private)
    ic(text)


if __name__ == '__main__':
    main()
