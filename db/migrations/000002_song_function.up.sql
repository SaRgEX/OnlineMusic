CREATE OR REPLACE FUNCTION song_filter(
    _name text,
    _performer int,
    _lyric text,
    _link text,
    _start_date date,
    _end_date date,
    _cursor int,
    _page_size int
)
    RETURNS TABLE
            (
                id           int,
                name         text,
                release_date date,
                lyric        text,
                link         text,
                performer      text
            )
AS
$$
BEGIN
    RETURN QUERY
        SELECT s.id, s.name, s.release_date, s.lyric, s.link, p.name
        FROM song AS s
        JOIN music_performer AS p ON p.id = s.performer_id
        WHERE (COALESCE(_name, '') = '' OR s.name ILIKE ('%' || _name || '%'))
          AND s.release_date BETWEEN COALESCE(_start_date, '1900-01-01') AND COALESCE(_end_date, '9999-12-31')
          AND (COALESCE(_lyric, '') = '' OR s.lyric ILIKE ('%' || _lyric || '%'))
          AND (COALESCE(_link, '') = '' OR s.link ILIKE ('%' || _link || '%'))
          AND s.id > _cursor
          AND (p.id = _performer OR _performer = -1)
        LIMIT _page_size;
END;
$$ LANGUAGE plpgsql;