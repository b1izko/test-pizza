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
const bolt_1 = require("@slack/bolt");
const bot_1 = require("./bot");
const bot = new bot_1.Bot();
const app = new bolt_1.App({
    appToken: bot.config.app_token,
    signingSecret: bot.config.sigint_secret,
    socketMode: bot.config.socketMode,
    token: bot.config.bot_token,
});
// Listen for users opening your App Home
app.event('im_open', ({ say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        //await ack();
        yield say("⚠ HELP ⚠\nUse the commands to enter an order: /<command> <parameter>\n" +
            "Example: /name Pepperoni");
        yield say("⚠ List of commands ⚠\n" +
            "/name - to enter the name of the pizza\n" +
            "/size - to enter the size of the pizza\n" +
            "/dough -  to enter the dough of the pizza\n" +
            "/extra -  to enter the name of any extra toppings\n" +
            "/address - to enter the delivery address\n" +
            "/comment - to enter an order comment\n" +
            "/accept - to confirm your order");
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
//app.message
app.command("/help", ({ ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        yield say("⚠ HELP ⚠\nUse the commands to enter an order: /<command> <parameter>\n" +
            "Example: /name Pepperoni");
        yield say("⚠ List of commands ⚠\n" +
            "/name - to enter the name of the pizza\n" +
            "/size - to enter the size of the pizza\n" +
            "/dough -  to enter the dough of the pizza\n" +
            "/extra -  to enter the name of any extra toppings\n" +
            "/address - to enter the delivery address\n" +
            "/comment - to enter an order comment\n" +
            "/accept - to confirm your order");
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/name", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (command.text == "") {
            yield say("❌ Error! Enter the name of the pizza\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Have you chosen: " + txt);
            bot.Set(command.user_id, txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                if (order.comment == "") {
                    yield say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅");
                }
                else {
                    yield say("☛ Your comment to order: " + order.comment);
                }
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/size", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (command.text == "") {
            yield say("❌ Error! Enter the size of the pizza\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Have you chosen: " + txt);
            bot.Set(command.user_id, undefined, txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                if (order.comment == "") {
                    yield say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅");
                }
                else {
                    yield say("☛ Your comment to order: " + order.comment);
                }
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/dough", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (txt.length == 0) {
            yield say("❌ Error! Enter the pizza dough\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Have you chosen: " + txt);
            bot.Set(command.user_id, undefined, undefined, txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                if (order.comment == "") {
                    yield say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅");
                }
                else {
                    yield say("☛ Your comment to order: " + order.comment);
                }
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/extra", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (txt.length == 0) {
            yield say("❌ Error! Enter the name of any extra toppings\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Have you chosen: " + txt);
            bot.Set(command.user_id, undefined, undefined, undefined, txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                if (order.comment == "") {
                    yield say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅");
                }
                else {
                    yield say("☛ Your comment to order: " + order.comment);
                }
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/address", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (txt.length == 0) {
            say("❌ Error! Enter the delivery address\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Your address: " + txt);
            bot.Set(command.user_id, undefined, undefined, undefined, txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                if (order.comment == "") {
                    yield say("🤔 Maybe you need a comment on the order? 🤔\n ✅ Use /comment - to enter an order comment! ✅");
                }
                else {
                    yield say("☛ Your comment to order: " + order.comment);
                }
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/comment", ({ command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        let txt = command.text;
        if (txt.length == 0) {
            say("❌ Error! Wrong comment\n ❗ For more information, use /help");
        }
        else {
            yield say("✅ Your comment: " + txt);
            if (bot.IsOrderComplete(command.user_id)) {
                let order = bot.Get(command.user_id);
                yield say("🍕 Your order 🍕" +
                    "\nPizza name: " + order.name +
                    "\nPizza size: " + order.size +
                    "\nPizza dough: " + order.dough +
                    "\nPizza extra toppings: " + order.extra +
                    "\nDelivery address: " + order.address);
                yield say("☛ Your comment to order: " + order.comment);
                yield say("Check your order and use /accept - to confirm your order");
            }
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
app.command("/accept", ({ client, command, ack, say }) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield ack();
        if (!bot.IsOrderComplete(command.user_id)) {
            yield say("❌ Error! Complete your order to confirm\n ❗ For more information, use /help");
        }
        else {
            bot.SendOrder(command.user_id);
            yield say("✅ Your order has been sent to the manager!");
        }
    }
    catch (error) {
        console.log("err");
        console.error(error);
    }
}));
(() => __awaiter(void 0, void 0, void 0, function* () {
    // Start bot app
    yield app.start();
    console.log('⚡️ Pizza app is running!');
}))();
