DROP TABLE IF EXISTS "played_cards";
CREATE TABLE "public"."played_cards" (
    "card_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "turn_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    CONSTRAINT "played_cards_pkey" PRIMARY KEY ("card_id", "turn_id")
) WITH (oids = false);

INSERT INTO "played_cards" ("card_id", "turn_id") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'c4fdc129-ed48-48df-a262-6be92a3acb12'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'1565f159-b6a2-4108-bc29-456dd05a8ac4'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'02f5e742-0221-484d-a796-c74f4691e693');

DROP TABLE IF EXISTS "turns";
CREATE TABLE "public"."turns" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "status" character varying(20) NOT NULL,
    "creation_date" timestamp NOT NULL,
    "formation_date" timestamp,
    "completion_date" timestamp,
    "moderator_id" uuid,
    "customer_id" uuid NOT NULL,
    "phase" character varying(50) NOT NULL,
    CONSTRAINT "turns_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "turns" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "phase") VALUES
('c4fdc129-ed48-48df-a262-6be92a3acb12',	'Введен',	'2023-10-09',	'2023-11-09',	'2023-12-09',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Развитие'),
('3c0a95f9-2e96-421b-903d-3bb9262cd77e',	'В работе',	'2023-09-24',	'2023-09-27',	'2023-09-28',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Кормовая база'),
('02f5e742-0221-484d-a796-c74f4691e693',	'Завершен',	'2023-10-09',	'2023-10-15',	'2023-10-20',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Питание'),
('1565f159-b6a2-4108-bc29-456dd05a8ac4',	'Отменен',	'2023-09-15',	'2023-09-15',	'2023-09-16',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Вымирание'),
('03d54ad2-0b82-4d49-9ca0-e67add804f4d',	'Удален',	'2023-10-08',	'2023-10-09',	'2023-10-10',	'87d54d58-1e24-4cca-9c83-bd2523902729',	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'Развитие'),
('6899a458-f86f-4cad-95ee-43c0e646d0ec',	'черновик',	'2023-10-29',	NULL,	NULL,	NULL,	'2d217868-ab6d-41fe-9b34-7809083a2e8a',	'');

DROP TABLE IF EXISTS "cards";
CREATE TABLE "public"."cards" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "name" character varying(50) NOT NULL,
    "image_url" character varying(100) NOT NULL,
    "type" character varying(50) NOT NULL,
    "description" character varying(200) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    CONSTRAINT "cards_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "cards" ("uuid", "name", "image_url", "type", "description", "is_deleted") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca', 'Большой',	'http://localhost:8080/image/img1.png',	'Хищник', 'Данное животное может быть съедено только большим хищником', 't'),
('b9778018-9c13-46fd-b785-4a803dc8be0b', 'Взаимодействие',	'http://localhost:8080/image/img2.png',	'Хищник', 'Сыграть одновременно на пару существ. Когда одно получит еду, то другое получит вне очереди', 'f'),
('18ab9f76-7648-49d2-857d-75ffddf13bea', 'Острое зрение',	'http://localhost:8080/image/img4.png',	'Жировой запас', 'Хищник имеющий этой свойство может съесть животной со свойством камуфляж', 'f');

DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "login" character varying(30) NOT NULL,
    "password" character varying(30) NOT NULL,
    "name" character varying(50) NOT NULL,
    "moderator" boolean NOT NULL,
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "login", "password", "name", "moderator") VALUES
('2d217868-ab6d-41fe-9b34-7809083a2e8a',	'user1',	'user1',	'Игрок',	'f'),
('87d54d58-1e24-4cca-9c83-bd2523902729',	'user2',	'user2',	'Модератор',	't');

ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_turn" FOREIGN KEY (turn_id) REFERENCES turns(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_card" FOREIGN KEY (card_id) REFERENCES cards(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;