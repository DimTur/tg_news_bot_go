┌─────────────────────────┐
│     Administrator       │
│       (Telegram)        │
└───────────┬─────────────┘
            │
            │
            │
            ▼
┌─────────────────────────┐
│      Telegram Bot       │
└───────────┬─────────────┘
            │
            │
            │
            ▼
┌───────────┼─────────────┐
│           │             │
│           │             │
│           │             │
│   ┌───────▼───────┐     │
│   │   Notifier    │     │
│   │(API ChatGPT)  │     │
│   └───────┬───────┘     │
│           │             │
│           │             │
│   ┌───────▼───────┐     │
│   │   Fetcher     │     │
│   │   (RSS)       │     │
│   └───────┬───────┘     │
│           │             │
│           │             │
│   ┌───────▼─────────────▼───┐
│   │    Database (PostgreSQL)│
│   │     ┌──────────────┐    │
│   │     │   articles   │    │
│   │     │   sources    │    │
│   │     └──────────────┘    │
│   └─────────────────────────┘
└─────────────────────────────┘