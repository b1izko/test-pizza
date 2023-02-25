"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bot = void 0;
const amqplib_1 = require("amqplib");
const config_1 = require("./config");
class Pizza {
    constructor(_name = "", _size = "", _dough = "", _extra = "", _address = "", _comment = "") {
        this.name = _name;
        this.size = _size;
        this.dough = _dough;
        this.extra = _extra;
        this.address = _address;
        this.comment = _comment;
    }
}
class Bot {
    constructor() {
        this.config = new config_1.Config();
        //this.config.Load();
        this.orders = new Map();
    }
    // IsOrderComplete checks is order filled by the user 
    IsOrderComplete(user) {
        let pizza = this.orders.get(user);
        if (typeof (pizza === null || pizza === void 0 ? void 0 : pizza.name) == "string" && (pizza === null || pizza === void 0 ? void 0 : pizza.name) != "" &&
            typeof (pizza === null || pizza === void 0 ? void 0 : pizza.size) == "string" && (pizza === null || pizza === void 0 ? void 0 : pizza.size) != "" &&
            typeof (pizza === null || pizza === void 0 ? void 0 : pizza.dough) == "string" && (pizza === null || pizza === void 0 ? void 0 : pizza.dough) != "" &&
            typeof (pizza === null || pizza === void 0 ? void 0 : pizza.extra) == "string" && (pizza === null || pizza === void 0 ? void 0 : pizza.extra) != "" &&
            typeof (pizza === null || pizza === void 0 ? void 0 : pizza.address) == "string" && (pizza === null || pizza === void 0 ? void 0 : pizza.address) != "") {
            return true;
        }
        return false;
    }
    // Set is setter for setting user's order
    Set(user, _name = "", _size = "", _dough = "", _extra = "", _address = "", _comment = "") {
        let order = this.Get(user);
        let newOrder = new Pizza(_name == "" ? order.name : _name, _size == "" ? order.size : _size, _dough == "" ? order.dough : _dough, _extra == "" ? order.extra : _extra, _address == "" ? order.address : _address, _comment == "" ? order.comment : _comment);
        this.orders.set(user, newOrder);
    }
    // Get is getter for getting user's order 
    Get(user) {
        let order = this.orders.get(user);
        return new Pizza((order === null || order === void 0 ? void 0 : order.name) === undefined ? "" : order === null || order === void 0 ? void 0 : order.name, (order === null || order === void 0 ? void 0 : order.size) === undefined ? "" : order === null || order === void 0 ? void 0 : order.size, (order === null || order === void 0 ? void 0 : order.dough) === undefined ? "" : order === null || order === void 0 ? void 0 : order.dough, (order === null || order === void 0 ? void 0 : order.extra) === undefined ? "" : order === null || order === void 0 ? void 0 : order.extra, (order === null || order === void 0 ? void 0 : order.address) === undefined ? "" : order === null || order === void 0 ? void 0 : order.address, (order === null || order === void 0 ? void 0 : order.comment) === undefined ? "" : order === null || order === void 0 ? void 0 : order.comment);
    }
    // SendOrder sends the order to RabbitMQ
    SendOrder(user) {
        return __awaiter(this, void 0, void 0, function* () {
            const connection = yield (0, amqplib_1.connect)(this.config.rabbitmq.url);
            // Create a channel
            const channel = yield connection.createChannel();
            // Makes the queue available to the client
            yield channel.assertQueue(this.config.rabbitmq.queue);
            // Send some messages to the queue
            if (this.IsOrderComplete(user)) {
                return;
            }
            let order = this.Get(user);
            let msg = {
                user: user,
                name: order.name,
                size: order.size,
                dough: order.dough,
                extra: order.extra,
                address: order.address,
                comment: order.comment
            };
            channel.sendToQueue(this.config.rabbitmq.queue, Buffer.from(JSON.stringify(msg)));
            // Init the consumer
            let consumer = (channel) => (msg) => {
                if (msg) {
                    // Display the received message
                    console.log(msg.content.toString());
                    // Acknowledge the message
                    channel.ack(msg);
                }
            };
            // Start the consumer
            yield channel.consume(this.config.rabbitmq.queue, consumer(channel));
        });
    }
}
exports.Bot = Bot;
