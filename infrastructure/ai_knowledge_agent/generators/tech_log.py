"""
Tech log data generator.
Creates structured log data using Faker.
"""

import random
from datetime import timedelta
from faker import Faker

from data_models import LogData, LogEntry

fake = Faker('de_DE')

# Log message templates
LOG_TEMPLATES = [
    ("INFO", "Service started successfully"),
    ("INFO", "Health check passed"),
    ("INFO", "Request processed in {ms}ms"),
    ("WARN", "High memory usage detected: {pct}%"),
    ("WARN", "Slow query detected: {sec} seconds"),
    ("ERROR", "Connection timeout to database"),
    ("ERROR", "Failed to write to disk: {err}"),
    ("DEBUG", "Cache hit rate: {pct}%"),
    ("DEBUG", "Active connections: {cnt}"),
]


def _format_message(message: str) -> str:
    """Fill in placeholders in log message templates."""
    if '{ms}' in message:
        return message.format(ms=random.randint(10, 5000))
    elif '{pct}' in message:
        return message.format(pct=random.randint(50, 95))
    elif '{sec}' in message:
        return message.format(sec=round(random.uniform(1.0, 10.0), 2))
    elif '{cnt}' in message:
        return message.format(cnt=random.randint(10, 500))
    elif '{err}' in message:
        return message.format(err="unknown error")
    return message


def create_log_data() -> LogData:
    """
    Generate structured tech log data.
    
    Returns:
        LogData: Complete log with entries populated.
    """
    server_prefix = random.choice(['web', 'db', 'api', 'cache', 'worker'])
    server_name = f"{server_prefix}-{random.randint(1, 10)}"
    log_date = fake.date_time_between(start_date='-30d', end_date='now')

    num_entries = random.randint(10, 30)
    entries = []
    current_time = log_date

    for _ in range(num_entries):
        level, message_template = random.choice(LOG_TEMPLATES)
        message = _format_message(message_template)

        entries.append(LogEntry(
            timestamp=current_time,
            level=level,
            server=server_name,
            message=message,
        ))
        current_time += timedelta(seconds=random.randint(1, 300))

    return LogData(
        server_name=server_name,
        log_date=log_date,
        entries=entries,
        uptime_hours=random.randint(1, 720),
    )
