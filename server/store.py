from abc import ABCMeta, abstractmethod


class Store(object):
	__metaclass__ = ABCMeta

	@abstractmethod
	def get(self, key): pass

	@abstractmethod
	def put(self, key, value): pass


class InMemoryStore(Store):

	def __init__(self):
		self.db = {}

	def get(self, key):
		return self.db[key]

	def put(self, key, value):
		self.db[key] = value