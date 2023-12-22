DROP TABLE IF EXISTS "played_cards";
CREATE TABLE "public"."played_cards" (
    "card_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    "turn_id" uuid DEFAULT gen_random_uuid() NOT NULL,
    CONSTRAINT "played_cards_pkey" PRIMARY KEY ("card_id", "turn_id")
) WITH (oids = false);

INSERT INTO "played_cards" ("card_id", "turn_id") VALUES
('f269cc75-9d07-4ded-bbb4-02536f742220',	'ed170485-5acf-4779-9e88-786664a6d23e'),
('c502737b-ec5c-41f6-b656-b8b08caedfc4',	'19199780-b53b-426e-85a2-281b2cb634a5'),
('320877f2-702c-4c8a-8d47-f17bbc51af55',	'7db3ffb5-1523-4813-b911-f19fec54e218'),
('f269cc75-9d07-4ded-bbb4-02536f742220',	'a3246620-e671-4602-871e-9ff0c84ec8be');

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
    "take_food" bigint NOT NULL,
    CONSTRAINT "turns_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "turns" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "phase", "sending_status") VALUES
('ed170485-5acf-4779-9e88-786664a6d23e',	'сформирован',	'2023-12-22 20:50:44.309105',	NULL,	NULL,	'e50e51a7-d4df-4136-89b2-79d43a577721',	'604bab40-0f84-43a1-87b4-169472eb5ce0',	NULL,	NULL),
('19199780-b53b-426e-85a2-281b2cb634a5',	'сформирован',	'2023-12-22 20:51:40.225756',	NULL,	NULL,	'e50e51a7-d4df-4136-89b2-79d43a577721',	'3fdb245e-d7ea-4ea0-915e-d3eafb8b6b7e',	NULL,	NULL),
('a3246620-e671-4602-871e-9ff0c84ec8be',	'сформирован',	'2023-12-22 20:52:34.068756',	NULL,	NULL,	'e50e51a7-d4df-4136-89b2-79d43a577721',	'3fdb245e-d7ea-4ea0-915e-d3eafb8b6b7e',	NULL,	NULL),
('7db3ffb5-1523-4813-b911-f19fec54e218',	'сформирован',	'2023-12-22 20:52:46.970346',	NULL,	NULL,	'e50e51a7-d4df-4136-89b2-79d43a577721',	'3fdb245e-d7ea-4ea0-915e-d3eafb8b6b7e',	NULL,	NULL);

DROP TABLE IF EXISTS "cards";
CREATE TABLE "public"."cards" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "name" character varying(50) NOT NULL,
    "image_url" character varying(100) NOT NULL,
    "type" character varying(50) NOT NULL,
    "description" character varying(200) NOT NULL,
    "is_deleted" boolean DEFAULT false NOT NULL,
    "need_food" bigint NOT NULL,
    CONSTRAINT "cards_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "cards" ("uuid", "image_url", "is_deleted", "type", "name", "description", "need_food") VALUES
('f269cc75-9d07-4ded-bbb4-02536f742220',	'http://localhost:8080/image/img9.png',	'f',	'Хищник',	'Большой',	'Данное животное может быть съедено только большим хищником',	1),
('c502737b-ec5c-41f6-b656-b8b08caedfc4',	'http://localhost:8080/image/img12.png',	'f',	'Не хищник',	'Большой',	'Сыграть одновременно на пару существ. Когда одно получит еду, то другое получит вге очереди',	2),
('320877f2-702c-4c8a-8d47-f17bbc51af55',	'http://localhost:8080/image/img6.png',	'f',	'Хищник',	'Взаимодействие',	'Сыграть одновременно на пару существ. Когда одно получит еду, то другое получит вне очереди',	1);

DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    "login" character varying(30) NOT NULL,
    "password" character varying(30) NOT NULL,
    "role" bigint,
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "login", "password", "role") VALUES
('604bab40-0f84-43a1-87b4-169472eb5ce0',	'antiho',	'Aa12345678',	1),
('e50e51a7-d4df-4136-89b2-79d43a577721',	'admin',	'Aa12345678',	2),
('3fdb245e-d7ea-4ea0-915e-d3eafb8b6b7e',	'string',	'string',	1);

ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_turn" FOREIGN KEY (turn_id) REFERENCES turns(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_card" FOREIGN KEY (card_id) REFERENCES cards(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;


INSERT INTO "notifications" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "notification_type") VALUES
('c4fdc129-ed48-48df-a262-6be92a3acb12',	'Введен',	'2023-10-09 00:00:00',	'2023-11-09 00:00:00',	'2023-12-09 00:00:00',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	'Срочное сообщение'),
('3c0a95f9-2e96-421b-903d-3bb9262cd77e',	'В работе',	'2023-09-24 00:00:00',	'2023-09-27 00:00:00',	'2023-09-28 00:00:00',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	'Еженедельное уведомление'),
('a10adc9b-4957-4769-a029-71b76c764e49',	'удалён',	'2023-10-29 00:00:00',	NULL,	NULL,	NULL,	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	''),
('1565f159-b6a2-4108-bc29-456dd05a8ac4',	'отклонён',	'2023-09-15 00:00:00',	'2023-09-15 00:00:00',	'2023-09-16 00:00:00',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	'Уведомление о задолжности'),
('02f5e742-0221-484d-a796-c74f4691e693',	'завершён',	'2023-10-09 00:00:00',	'2023-10-15 00:00:00',	'2023-10-20 00:00:00',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	'Электронное напоминание'),
('200b2366-36d5-49b2-9770-85c1628c20f0',	'сформирован',	'2023-10-30 00:00:00',	'2023-10-30 00:00:00',	NULL,	NULL,	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	''),
('da900b13-b937-4f05-baf2-0c114dcb1967',	'завершён',	'2023-12-06 05:58:56.052583',	'2023-12-06 05:59:16.93331',	'2023-12-06 06:01:14.437566',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	''),
('59128f4e-4011-4bef-9a6d-13ab40659d27',	'завершён',	'2023-12-05 23:20:03.994931',	'2023-12-05 23:24:19.679904',	'2023-12-06 06:43:18.778483',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	''),
('4838989d-3d71-4d9e-a8d9-7f12dd76592d',	'завершён',	'2023-12-06 06:01:25.483304',	'2023-12-06 06:01:35.952638',	'2023-12-06 07:24:32.899471',	'3e9e10cc-a591-4d6b-a017-ac23eeffee39',	'967db760-fef5-4eb1-b26a-dc25c3cd40a9',	'');



INSERT INTO "recipients" ("uuid", "fio", "image_url", "email", "age", "adress", "is_deleted") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'Олег Орлов Никитович',	'localhost:9000/images/men1.jpg',	'OlegO@mail.ru',	27,	'Москва, ул. Измайловская, д.13, кв.54',	't'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'Василий Гречко Валентинович',	'localhost:9000/images/men2.jpg',	'Grechko_101@mail.ru',	31,	'Москва, ул. Тверская, д.25, кв.145',	'f'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'Александр Лейко Кириллович',	'localhost:9000/images/men3.jpg',	'Alek221@mail.ru',	37,	'Москва, ул. Изюмская, д.15, кв.89',	'f'),
('365f46c8-b498-47b9-92d3-97319ff62711',	'Андрей Отрис Даниллович',	'localhost:9000/images/men1.jpg',	'Andr1@mail.ru',	32,	'Москва, ул. Изюмская, д.15, кв.79',	'f'),
('9b8914a6-c599-450d-893d-b8ebb766dd07',	'Кирилл Лейка Кириллович',	'localhost:9000/images/men4.jpg',	'KriLeik@mail.ru',	30,	'Москва, ул. Бутовская, д.15, кв.79',	'f'),
('b06a0b5a-6ede-4636-a97b-83976dc10575',	'Андрей Ермолин Данилович',	'localhost:9000/images/men5.jpg',	'andrErm@gmail.com',	33,	'Москва, ул. Клинская, д.22, кв.12',	'f');

INSERT INTO "notification_contents" ("recipient_id", "notification_id") VALUES
('4bea0842-bcb8-416e-9a63-d89a63e978ca',	'c4fdc129-ed48-48df-a262-6be92a3acb12'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'1565f159-b6a2-4108-bc29-456dd05a8ac4'),
('18ab9f76-7648-49d2-857d-75ffddf13bea',	'02f5e742-0221-484d-a796-c74f4691e693'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'200b2366-36d5-49b2-9770-85c1628c20f0'),
('b9778018-9c13-46fd-b785-4a803dc8be0b',	'59128f4e-4011-4bef-9a6d-13ab40659d27'),
('b06a0b5a-6ede-4636-a97b-83976dc10575',	'da900b13-b937-4f05-baf2-0c114dcb1967'),
('b06a0b5a-6ede-4636-a97b-83976dc10575',	'4838989d-3d71-4d9e-a8d9-7f12dd76592d');

