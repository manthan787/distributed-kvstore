from abc import ABCMeta, abstractmethod


class Store(object):
	__metaclass__ = ABCMeta

	@abstractmethod
	def get(self, key): pass

	@abstractmethod
	def put(self, key, value): pass

	@abstractmethod
	def batch_put(self, kvs): pass


class InMemoryStore(Store):

	def __init__(self):
		self.db = {}

	def get(self, key):
		return self.db[key]

	def put(self, key, value):
		if key in db: return False
		self.db[key] = value
		return True

	def batch_put(self, kvs):
		stats = {'keys_added': 0, 'keys_failed': 0}
		for kv in kvs:
			if self.put(kv["key"], kv["value"]): stats[keys_added] += 1
			else: stats[keys_failed] += 1
		return stats