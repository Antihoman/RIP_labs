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
    "phase" character varying(50),
    "take_food" bigint NOT NULL,
    "sending_status" character varying(40),
    CONSTRAINT "turns_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "turns" ("uuid", "status", "creation_date", "formation_date", "completion_date", "moderator_id", "customer_id", "phase", "sending_status", "take_food") VALUES
('ed170485-5acf-4779-9e88-786664a6d23e',	'сформирован',	'2023-12-22 20:50:44.309105',	'2023-12-22 20:54:44.309105',	NULL,	NULL,	'de327707-cc73-4a0d-adc7-2ae760b9e648',	NULL,	NULL, 1),
('19199780-b53b-426e-85a2-281b2cb634a5',	'сформирован',	'2023-12-22 20:51:40.225756',	'2023-12-22 20:56:44.309105',	NULL,	NULL,	'de327707-cc73-4a0d-adc7-2ae760b9e648',	NULL,	NULL, 2),
('a3246620-e671-4602-871e-9ff0c84ec8be',	'сформирован',	'2023-12-22 20:52:34.068756',	NULL,	NULL,	NULL,	'702a4657-4ee6-49f2-8bc4-bcb2d61f3e5b',	NULL,	NULL, 0),
('7db3ffb5-1523-4813-b911-f19fec54e218',	'сформирован',	'2023-12-22 20:52:46.970346',	NULL,	NULL,	NULL,	'702a4657-4ee6-49f2-8bc4-bcb2d61f3e5b',	NULL,	NULL, 1);

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
    "password" character varying(40) NOT NULL,
    "role" bigint,
    CONSTRAINT "users_pkey" PRIMARY KEY ("uuid")
) WITH (oids = false);

INSERT INTO "users" ("uuid", "login", "password", "role") VALUES
('de327707-cc73-4a0d-adc7-2ae760b9e648',	'antiho',	'4fad0bc2df9b917b93047316f9f852eb8760aea3',	1),
('702a4657-4ee6-49f2-8bc4-bcb2d61f3e5b',	'user',	'4fad0bc2df9b917b93047316f9f852eb8760aea3',	1),
('c4c0ea17-934d-4882-8f6d-a57706088d8d',	'admin',	'4fad0bc2df9b917b93047316f9f852eb8760aea3',	2);

ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_turn" FOREIGN KEY (turn_id) REFERENCES turns(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."played_cards" ADD CONSTRAINT "fk_played_cards_card" FOREIGN KEY (card_id) REFERENCES cards(uuid) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_customer" FOREIGN KEY (customer_id) REFERENCES users(uuid) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."turns" ADD CONSTRAINT "fk_turns_moderator" FOREIGN KEY (moderator_id) REFERENCES users(uuid) NOT DEFERRABLE;