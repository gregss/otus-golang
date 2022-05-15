package storage

import "time"

type Event struct {
	ID          int
	Title       string
	Time        time.Time
	Duration    time.Duration
	Description string
	UserID      int
	NotifyTime  time.Time
}

/*ID - уникальный идентификатор события (можно воспользоваться UUID);
Заголовок - короткий текст;
Дата и время события;
Длительность события (или дата и время окончания);
Описание события - длинный текст, опционально;
ID пользователя, владельца события;
За сколько времени высылать уведомление, опционально.*/
