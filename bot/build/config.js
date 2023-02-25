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
exports.Config = void 0;
const js_yaml_1 = require("js-yaml");
const promises_1 = require("fs/promises");
class RabbitMQ {
    constructor(_url = "amqp://guest:guest@localhost:5672/", _queue = "go") {
        this.url = _url;
        this.queue = _queue;
    }
}
class Config {
    constructor(_app_token = "xapp-1-A04NQJV5JHM-4799882804385-35b612da99e179f0e2a39250200ad2b9a9aca147c1d14624bd70ad79423edfe3", _bot_token = "xoxb-4798470422000-4787063403475-QLvM5uJP3pnoZmNNYPUS0IRf", _sigint_secret = "1c3e64beb1dcf9147a83d2c3e9f44ae8", _rabbitmq_url) {
        this.app_token = _app_token;
        this.bot_token = _bot_token;
        this.sigint_secret = _sigint_secret;
        this.socketMode = true;
        this.rabbitmq = new RabbitMQ(_rabbitmq_url);
    }
    Load() {
        return __awaiter(this, void 0, void 0, function* () {
            let yaml = (0, js_yaml_1.load)(yield (0, promises_1.readFile)("config.yaml", "utf-8"));
            this.app_token = yaml.app_token;
            this.bot_token = yaml.bot_token;
            this.sigint_secret = yaml.sigint_secret;
            this.rabbitmq = yaml.rabbitmq;
            console.log("YAML: ", yaml);
            console.log("Readed config: ", this);
            return this;
        });
    }
}
exports.Config = Config;
