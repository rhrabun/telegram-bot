from telegram.ext import (
    CommandHandler,
    MessageHandler,
    Filters, 
    ConversationHandler
)

from ..states import SCRIPTS_STATE
from ..handlers import (
    default_handler,
    scripts_handler,
    run_lambda_handler
)


scripts_state_machine = ConversationHandler(
            entry_points = [
                CommandHandler(
                    'scripts',
                    callback=scripts_handler
                )
            ],
            states={
                SCRIPTS_STATE: [
                    MessageHandler(
                        callback=run_lambda_handler,
                        filters=(Filters.text & Filters.regex('^scripts-dev'))
                    )
                ]
            },
            fallbacks=[
                CommandHandler(
                    'cancel',
                    callback=default_handler
                )
            ]
        )
