import { App } from "@slack/bolt";
import client, {Connection, Channel, ConsumeMessage} from "amqplib";


//to .env!
const SLACK_APP_TOKEN = "";
const SLACK_BOT_TOKEN = "";
const SLACK_BOT_SIGINT_SECRET = "";

const app = new App({
  appToken: SLACK_APP_TOKEN,
  signingSecret: SLACK_BOT_SIGINT_SECRET,
  socketMode: true,
  token: SLACK_BOT_TOKEN,
});

type Pizza = {
    name: string;
    size: string;
    dough: string;
    extra: string;
    address: string;
    comment: string;
}

const pizzaOrders = new Map<string, Pizza>

function isOrderComplete(user: string): boolean{
    let buffer = pizzaOrders.get(user)
    if (
        typeof buffer?.name == "string" && buffer?.name != "" &&
        typeof buffer?.size == "string" && buffer?.size != "" &&
        typeof buffer?.dough == "string" && buffer?.dough != "" &&
        typeof buffer?.extra == "string" && buffer?.extra != "" &&
        typeof buffer?.address == "string" && buffer?.address != ""
    ) {
        return true
    }
    return false
}

const sendOrder = (user: string, channel: Channel): boolean => {
    if (!isOrderComplete(user)) {
        return false
    }
    
    let buffer = pizzaOrders.get(user)  
    let msg = { 
        user: user, 
        name: buffer?.name,
        size: buffer?.size,
        dough: buffer?.dough,
        extra: buffer?.extra,
        address: buffer?.address,
        comment: buffer?.comment 
    }

    channel.sendToQueue("go", Buffer.from(JSON.stringify(msg)))

    return true
  }

const consumer = (channel: Channel) => (msg: ConsumeMessage | null): void => {
    if (msg) {
      // Display the received message
      console.log(msg.content.toString())
      // Acknowledge the message
      channel.ack(msg)
    }
  }

app.command("/help", async ({ ack, say }) =>{
    try {
        await ack();
        await say("⚠ HELP ⚠\nUse the commands to enter an order: /<command> <parameter>\n" + 
                "Example: /name Pepperoni")
        await say("⚠ List of commands ⚠\n" + 
                "/name - to enter the name of the pizza\n" + 
                "/size - to enter the size of the pizza\n" + 
                "/dough -  to enter the dough of the pizza\n" + 
                "/extra -  to enter the name of any extra toppings\n" + 
                "/address - to enter the delivery address\n" + 
                "/comment - to enter an order comment\n" + 
                "/accept - to confirm your order")
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/name", async ({ command, ack, say }) =>{
    try {
        await ack();
        let txt = command.text
        if (command.text == "") {
            await say("❌ Error! Enter the name of the pizza\n ❗ For more information, use /help")
        } else {
            await say("✅ Have you chosen: " + txt)
            
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : txt,
                size : buffer?.size === undefined ? "" : buffer?.size,
                dough : buffer?.dough === undefined ? "" : buffer?.dough,
                extra : buffer?.extra === undefined ? "" : buffer?.extra,
                address : buffer?.address === undefined ? "" : buffer?.address,
                comment : buffer?.comment === undefined ? "" : buffer?.comment,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                if (buffer?.comment === undefined || buffer?.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + buffer?.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/size", async ({ command, ack, say }) => {
    try {
        await ack();
        let txt = command.text
        if (txt.length == 0) {
            await say("❌ Error! Enter the size of the pizza\n ❗ For more information, use /help")
        } else {
            await say("✅ Have you chosen: " + txt)
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : buffer?.name === undefined ? "" : buffer?.name,
                size : txt,
                dough : buffer?.dough === undefined ? "" : buffer?.dough,
                extra : buffer?.extra === undefined ? "" : buffer?.extra,
                address : buffer?.address === undefined ? "" : buffer?.address,
                comment : buffer?.comment === undefined ? "" : buffer?.comment,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                if (buffer?.comment === undefined || buffer?.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + buffer?.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/dough", async ({ command, ack, say }) => {
    try {
        await ack();
        let txt = command.text
        if (txt.length == 0) {
            await say("❌ Error! Enter the pizza dough\n ❗ For more information, use /help")
        } else {
            await say("✅ Have you chosen: " + txt)
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : buffer?.name === undefined ? "" : buffer?.name,
                size : buffer?.size === undefined ? "" : buffer?.size,
                dough : txt,
                extra : buffer?.extra === undefined ? "" : buffer?.extra,
                address : buffer?.address === undefined ? "" : buffer?.address,
                comment : buffer?.comment === undefined ? "" : buffer?.comment,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                if (buffer?.comment === undefined || buffer?.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + buffer?.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/extra", async ({ command, ack, say }) => {
    try {
        await ack();
        let txt = command.text
        if (txt.length == 0) {
            await say("❌ Error! Enter the name of any extra toppings\n ❗ For more information, use /help")
        } else {
            await say("✅ Have you chosen: " + txt)
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : buffer?.name === undefined ? "" : buffer?.name,
                size : buffer?.size === undefined ? "" : buffer?.size,
                dough : buffer?.dough === undefined ? "" : buffer?.dough,
                extra : txt,
                address : buffer?.address === undefined ? "" : buffer?.address,
                comment : buffer?.comment === undefined ? "" : buffer?.comment,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                if (buffer?.comment === undefined || buffer?.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + buffer?.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/address", async ({ command, ack, say }) => {
    try {
        await ack();
        let txt = command.text
        if (txt.length == 0) {
            say("❌ Error! Enter the delivery address\n ❗ For more information, use /help")
        } else {
            await say("✅ Your address: " + txt)
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : buffer?.name === undefined ? "" : buffer?.name,
                size : buffer?.size === undefined ? "" : buffer?.size,
                dough : buffer?.dough === undefined ? "" : buffer?.dough,
                extra : buffer?.extra === undefined ? "" : buffer?.extra,
                address : txt,
                comment : buffer?.comment === undefined ? "" : buffer?.comment,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                if (buffer?.comment === undefined || buffer?.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + buffer?.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/comment", async ({ command, ack, say }) => {
    try {
        await ack();
        let txt = command.text
        if (isOrderComplete(command.user_id)) {
            await say("✅ Your comment: " + txt)
            let buffer = pizzaOrders.get(command.user_id) 
            let newOrder : Pizza = {
                name : buffer?.name === undefined ? "" : buffer?.name,
                size : buffer?.size === undefined ? "" : buffer?.size,
                dough : buffer?.dough === undefined ? "" : buffer?.dough,
                extra : buffer?.extra === undefined ? "" : buffer?.extra,
                address : buffer?.address === undefined ? "" : buffer?.address,
                comment : txt,
            }  
            pizzaOrders.set(command.user_id, newOrder)
            if (isOrderComplete(command.user_id)) {
                buffer = pizzaOrders.get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + buffer?.name + 
                    "\nPizza size: " + buffer?.size +
                    "\nPizza dough: " + buffer?.dough +
                    "\nPizza extra toppings: " + buffer?.extra +
                    "\nDelivery address: " + buffer?.address)
                await say("☛ Your comment to order: " + buffer?.comment)
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/accept", async ({ command, ack, say }) => {
    try {
        await ack();
        if (!isOrderComplete(command.user_id)) {
            await say("❌ Error! Complete your order to confirm\n ❗ For more information, use /help")
        } else {
            const connection: Connection = await client.connect("amqp://localhost")
            // Create a channel
            const channel: Channel = await connection.createChannel()
            // Makes the queue available to the client
            await channel.assertQueue("go")
            // Send some messages to the queue
            sendOrder(command.user_id, channel)
            // Start the consumer
            await channel.consume("go", consumer(channel))
        }

    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.start().catch((error) => {
  console.error(error);
  process.exit(1);
});