from telegram.ext import (
    CommandHandler,
    MessageHandler,
    Filters, 
    ConversationHandler
)

from ..states import SERVER_STATE
from ..handlers import (
    default_handler,
    server_handler,
    status_handler
)


server_state_machine = ConversationHandler(
            entry_points = [
                CommandHandler(
                    'server',
                    callback=server_handler
                )
            ],
            states={
                SERVER_STATE: [
                    CommandHandler(
                        'status',
                        callback=status_handler
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
