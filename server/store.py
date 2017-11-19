from abc import ABCMeta, abstractmethod

class Store(object):
    """ This abstract class defines the contract for Store concrete implementations """
    __metaclass__ = ABCMeta

    @abstractmethod
    def get_all(self):
        """ Return all the key-value pairs in the store """
        pass

    @abstractmethod
    def get(self, key):
        """ Return value for specified `key` """
        pass

    @abstractmethod
    def put(self, key, value):
        """ Put given key-value pair inside the Store """
        pass

    @abstractmethod
    def batch_put(self, kvs):
        """ Put list of key-value pairs inside the Store """
        pass

    @abstractmethod
    def batch_get(self, keys):
        """ Get a list of key-value pairs from the Store """
        pass

class InMemoryStore(Store):
    """ In-Memory concrete implementation for `Store` """

    def __init__(self, valid_encodings):
        self.encodings = {}
        self.data = {}
        self.valid_encodings = valid_encodings

    def get_all(self):
        for k in self.data:
            data, success = self.get({"data": k})
            if success: yield data

    def get(self, key):
        try:
            key_data = key['data']
            return {"key": self._get_payload(key_data), \
                    "value": self._get_payload(self.data[key_data])}, True
        except Exception as e:
            print e
            return {}, False

    def put(self, key, value):
        try:
            if not self._is_valid_kv(key, value): return False
            self.encodings[key['data']] = key['encoding']
            self.encodings[value['data']] = value['encoding']
            self.data[key['data']] = value['data']
        except KeyError as e:
            print "Key not found: " + e
            return False
        return True

    def _is_valid_kv(self, key, value):
        return key['data'] in self.data and key['encoding'] in self.valid_encodings \
                and value['encoding'] in self.valid_encodings

    def batch_put(self, kvs):
        stats = {'keys_added': 0, 'keys_failed': []}
        for kv in kvs:
            if self.put(kv["key"], kv["value"]): stats["keys_added"] += 1
            else: stats["keys_failed"].append(kv.get("key", "Missing-Key-Field"))
        return stats

    def batch_get(self, keys):
        for key in keys:
            result, success = self.get(key)
            if success: yield result

    def _get_payload(self, k):
        return {"encoding": self.encodings[k], "data": k}