#!/bin/bash

sqlite3 mlitedb.db "DELETE FROM config;DELETE FROM payment_systems;DELETE FROM cf_lang;DELETE FROM install;DELETE FROM my_admin_messages;DELETE FROM my_cash_requests;DELETE FROM my_cf_funding;DELETE FROM my_comments;DELETE FROM my_commission;DELETE FROM my_complex_votes;DELETE FROM my_dc_transactions;DELETE FROM my_holidays;DELETE FROM my_keys;DELETE FROM my_new_users;DELETE FROM my_node_keys;DELETE FROM my_notifications;DELETE FROM my_promised_amount;DELETE FROM my_table;DELETE FROM my_tasks;delete from block_chain where id < 254000; VACUUM;"