import psycopg2


def connect():
    """ Connect to the PostgreSQL database server """
    conn = None
    try:
        # read connection parameters
        params = {"host": "ec2-107-22-122-106.compute-1.amazonaws.com",
                  "database": "d8ce7g1pdv8bao",
                  "user": "vdxeyeechvtzph",
                  "password": "bb994a8f08701050c5b0c0fdb33ef02525dacf7b0ccc6f3dc227be088cb8db7b"}

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
