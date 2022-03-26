.open db.db

drop table users;

create table users(
                       name text not null primary key,
                       is_admin integer not null,
                       password text not null
);

insert into users values("root", 1, "1111");
insert into users values("Говард Лавкрафт", 0, "нъярлатхотеп");
insert into users values("Амброз Бирс", 0, "12345");