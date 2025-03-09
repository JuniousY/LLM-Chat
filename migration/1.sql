-- conversations table
create table conversations
(
    id           int auto_increment primary key,
    app_id       int                      not null,
    name         varchar(255)             not null,
    introduction text null,
    model_name   varchar(50) null,
    status       varchar(255) default "0" not null,
    from_user_id varchar(255)             not null,
    is_deleted   tinyint(1) default 0 not null,
    created_at   datetime                 not null,
    updated_at   datetime                 not null
) comment 'llm对话表';

create index index_1
    on conversations (app_id, from_user_id);

-- messages table
create table messages
(
    id                        int auto_increment primary key,
    app_id                    int,
    model_provider            varchar(255),
    model_id                  varchar(255),
    override_model_configs    text,
    conversation_id           int,
    inputs                    json                         not null,
    query                     text                         not null,
    message                   json                         not null,
    message_tokens            integer          default 0   not null,
    message_unit_price        numeric(10, 4)               not null,
    answer                    text                         not null,
    answer_tokens             integer          default 0   not null,
    answer_unit_price         numeric(10, 4)               not null,
    provider_response_latency double precision default 0   not null,
    total_price               numeric(10, 7),
    currency                  varchar(255)                 not null,
    from_source               varchar(255)                 not null,
    from_user_id              varchar(255)                 not null,
    created_at                datetime                     not null,
    updated_at                datetime                     not null,
    status                    varchar(255)     default "0" not null,
    error                     text
) comment 'llm对话消息表';

create index index_conversation
    on messages (conversation_id);