# demo-producer

This is a demo-producer application.

To deploy this application, execute the following commands:

  1. Clone repo:

    ```
    $ git clone https://github.com/damianjaniszewski/go-env
    $ cd demo-producer
    ```

  2. Install Godep package manager (git required to complete):

    ```
    $ go get github.com/tools/godep

    ```

  3. Create Godep package manager files:

    ```
    $ godep save
    ```

  3. Deploy to HPE Stackato PaaS

    ```
    $ stackato push --stack cflinuxfs2 -n
    ```
