import psycopg2
import config


def connect():
    """ Connect to the PostgreSQL database server """
    try:
        # read connection parameters
        params = {"host": config.host,
                  "database": config.database,
                  "user": config.user,
                  "password": config.password}

        # connect to the PostgreSQL server
        print('Connecting to the PostgreSQL database...')
        conn = psycopg2.connect(**params)

        # create a cursor
        cur = conn.cursor()

        return cur, conn

        # close the communication with the PostgreSQL
        # cur.close()
    except (Exception, psycopg2.DatabaseError) as error:
        print(error)
    # finally:
    #     if conn is not None:
    #         conn.close()
    #         print('Database connection closed.')


if __name__ == '__main__':
    connect()
