import pytz
import logging
import datetime

from telegram import Update, ReplyKeyboardRemove
from telegram.ext import (
    Updater,
    CommandHandler,
    MessageHandler,
    Filters,
    CallbackContext
)

from .config import TG_BOT_TOKEN, TG_CHAT_ID
from .state_nachines import (
    scripts_state_machine,
    server_state_machine
)
from .handlers import (
    default_handler,
    morning_handler
)


class BotObj(object):
    def __init__(self):
        # Set up logging
        logging.basicConfig(level=logging.INFO)

        # Set up bot
        self.bot = Updater(TG_BOT_TOKEN)
        print(self.bot)
        self.dispatcher = self.bot.dispatcher
        print(self.dispatcher)

        self.jobs = self.bot.job_queue


    def set_handlers(self):
        # Register handler for commands 'start'
        self.dispatcher.add_handler(CommandHandler(
            ['start'],
            callback=default_handler
            ))
        
        self.dispatcher.add_handler(scripts_state_machine)
        self.dispatcher.add_handler(server_state_machine)

        # Register scheduled handler, which sends morning message
        morning_time = datetime.time(
            hour=10,
            minute=00,
            second=00
        ).replace(tzinfo=pytz.timezone('Europe/Warsaw'))

        self.jobs.run_daily(
            morning_handler,
            time=morning_time
        )


    def main(self):
        self.set_handlers()
        self.bot.start_polling()
        self.bot.idle()
