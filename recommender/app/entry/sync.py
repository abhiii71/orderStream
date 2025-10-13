from kafka import KafkaConsumer 
import json 
import requests
from app.db.session import ReplicaSession
from app.db.models import Product, Interaction
from config.settings import PRODUCT_API, KAFKA_SERVER

def sync_products():
    consumer = KafkaConsumer("product_events", bootstrap_servers=KAFKA_SERVER)
    for message in consumer:
        event = json.loads(message.value)
        with ReplicaSession() as session:
            if event["type"] in ["product_created", "product_updated"]:
                product_data = event["data"]
                print(f"Product data: {product_data}") 