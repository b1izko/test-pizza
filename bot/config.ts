import { load } from "js-yaml";
import { readFile } from "fs/promises";

class RabbitMQ {
    url: string;
    queue: string;

    constructor (_url: string = "amqp://guest:guest@localhost:5672/", _queue: string = "go") {
        this.url = _url;
        this.queue = _queue;
    }
}

export class Config {
    // bot config
    app_token: string;
    bot_token: string;
    socketMode: boolean;
    sigint_secret: string;

    //rabbitmq config
    rabbitmq: RabbitMQ;

    constructor(
        _app_token: string = "xapp-1-A04NQJV5JHM-4799882804385-35b612da99e179f0e2a39250200ad2b9a9aca147c1d14624bd70ad79423edfe3",
        _bot_token: string = "xoxb-4798470422000-4787063403475-QLvM5uJP3pnoZmNNYPUS0IRf",
        _sigint_secret: string = "1c3e64beb1dcf9147a83d2c3e9f44ae8",
        _rabbitmq_url?: string,
    ) {
        this.app_token = _app_token;
        this.bot_token = _bot_token;
        this.sigint_secret = _sigint_secret;
        this.socketMode = true;
        this.rabbitmq = new RabbitMQ(_rabbitmq_url)
    }

    async Load() {
        let yaml = load(await readFile("config.yaml", "utf-8")) as Config
        this.app_token = yaml.app_token;
        this.bot_token = yaml.bot_token;
        this.sigint_secret = yaml.sigint_secret;
        this.rabbitmq = yaml.rabbitmq;

        console.log("YAML: ", yaml)
        console.log("Readed config: ", this)
        return this
    }
}

