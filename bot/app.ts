import { App } from "@slack/bolt"
import { Bot } from "./bot"

const bot = new Bot()

const app = new App({
    appToken: bot.config.app_token,
    signingSecret: bot.config.sigint_secret,
    socketMode: bot.config.socketMode,
    token: bot.config.bot_token,
});

//
// DO FORK:
// https://github.com/slackapi/sample-message-menus-node/blob/master/lib/bot.js
//

// Listen for users opening your App Home
app.event('im_open', async ({ say }) => {
    try {
        //await ack();
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
//app.message

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
            
            bot.Set(command.user_id, txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                if (order.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + order.comment)
                }
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/size", async ({ command, ack, say }) =>{
    try {
        await ack();
        let txt = command.text
        if (command.text == "") {
            await say("❌ Error! Enter the size of the pizza\n ❗ For more information, use /help")
        } else {
            await say("✅ Have you chosen: " + txt)
            
            bot.Set(command.user_id, undefined, txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                if (order.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + order.comment)
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
            
            bot.Set(command.user_id, undefined, undefined, txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                if (order.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + order.comment)
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
            
            bot.Set(command.user_id, undefined, undefined, undefined, txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                if (order.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + order.comment)
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
            
            bot.Set(command.user_id, undefined, undefined, undefined, txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                if (order.comment == "") {
                    await say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅")
                } else {
                    await say("☛ Your comment to order: " + order.comment)
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
        if (txt.length == 0) {
            say("❌ Error! Wrong comment\n ❗ For more information, use /help")
        } else {
            await say("✅ Your comment: " + txt)
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id)
                await say("🍕 Your order 🍕" + 
                    "\nPizza name: " + order.name + 
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address)
                await say("☛ Your comment to order: " + order.comment)
                await say ("Check your order and use /accept - to confirm your order")
            }
        }
    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

app.command("/accept", async ({ client, command, ack, say }) => {
    try {
        await ack();
        if (!bot.IsOrderComplete(command.user_id)) {
            await say("❌ Error! Complete your order to confirm\n ❗ For more information, use /help")
        } else {
            bot.SendOrder(command.user_id)
            await say("✅ Your order has been sent to the manager!")
        }

    } catch (error) {
        console.log("err")
        console.error(error)
    }
});

(async () => {
    // Start bot app
    await app.start();
  
    console.log('⚡️ Pizza app is running!');
  })();
