import { Connection, Channel, ConsumeMessage, connect } from "amqplib";
import { Config } from "./config";

class Pizza {
    name: string;
    size: string;
    dough: string;
    extra: string;
    address: string;
    comment: string;

    constructor(
        _name: string = "",
        _size: string = "",
        _dough: string = "",
        _extra: string = "",
        _address: string = "",
        _comment: string = "",
    ) {
        this.name = _name;
        this.size = _size;
        this.dough = _dough;
        this.extra = _extra;
        this.address = _address;
        this.comment = _comment;
    }
}

export class Bot {
    // config
    config : Config;
    
    // orders info
    private orders: Map<string, Pizza>;

    constructor() {
        this.config = new Config();
        //this.config.Load();
        this.orders = new Map<string, Pizza>();      
    }

    // IsOrderComplete checks is order filled by the user 
    IsOrderComplete(user: string): boolean {
        let pizza = this.orders.get(user); 
        if (
            typeof pizza?.name == "string" && pizza?.name != "" &&
            typeof pizza?.size == "string" && pizza?.size != "" &&
            typeof pizza?.dough == "string" && pizza?.dough != "" &&
            typeof pizza?.extra == "string" && pizza?.extra != "" &&
            typeof pizza?.address == "string" && pizza?.address != ""
        ) {
            return true;
        }
        return false;
    }

    // Set is setter for setting user's order
    Set(
        user: string, 
        _name: string = "",
        _size: string = "",
        _dough: string = "",
        _extra: string = "",
        _address: string = "",
        _comment: string = ""
        ) {
        let order = this.Get(user);
        let newOrder = new Pizza(
            _name == "" ? order.name : _name,  
            _size == "" ? order.size : _size,
            _dough == "" ? order.dough : _dough,
            _extra == "" ? order.extra : _extra,
            _address == "" ? order.address : _address,
            _comment == "" ? order.comment : _comment,
        )
        this.orders.set(user, newOrder)
    }

    // Get is getter for getting user's order 
    Get(user: string): Pizza {
        let order = this.orders.get(user)
        return new Pizza(
            order?.name === undefined ? "" : order?.name,
            order?.size === undefined ? "" : order?.size,
            order?.dough === undefined ? "" : order?.dough,
            order?.extra === undefined ? "" : order?.extra,
            order?.address === undefined ? "" : order?.address,
            order?.comment === undefined ? "" : order?.comment
        )
    }

    // SendOrder sends the order to RabbitMQ
    async SendOrder(user: string) {
        const connection: Connection = await connect(this.config.rabbitmq.url)
        // Create a channel
        const channel: Channel = await connection.createChannel()
        // Makes the queue available to the client
        await channel.assertQueue(this.config.rabbitmq.queue)
        // Send some messages to the queue
        if (this.IsOrderComplete(user)) {
            return 
        }

        let order = this.Get(user)
        let msg = { 
            user: user, 
            name: order.name,
            size: order.size,
            dough: order.dough,
            extra: order.extra,
            address: order.address,
            comment: order.comment 
        }

        channel.sendToQueue(this.config.rabbitmq.queue, Buffer.from(JSON.stringify(msg)))
        
        // Init the consumer
        let consumer = (channel: Channel) => (msg: ConsumeMessage | null): void => {
            if (msg) {
              // Display the received message
              console.log(msg.content.toString())
              // Acknowledge the message
              channel.ack(msg)
            }
        }
        // Start the consumer
        await channel.consume(this.config.rabbitmq.queue, consumer(channel))
    }
}