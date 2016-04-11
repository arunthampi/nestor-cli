## Nestor Power Development Toolkit

The Nestor CLI lets you create, debug and deploy Nestor Bot Powers.

## What is a Bot Power?

A Bot Power is functionality that you would like to add to
[Nestor](https://www.asknestor.me). Examples of powers include
[Github](https://www.asknestor.me/powers/github) and
[Trello](https://www.asknestor.me/powers/trello). But what if you wanted
to add your own?

Nestor CLI comes to the rescue.

## Installation

The Nestor CLI can be downloaded from
[here](https://dl.equinox.io/zerobotlabs/nestor/stable). Downloads are
available for OS X and Linux. (Windows coming soon).

Unpackage the archive and move the `nestor` binary to a well-known
`$PATH` (such as `/usr/local/bin`).

## Start creating your own power

To start creating your own power, run `nestor new <Power Name>`. So if
you wanted to call your power "Hello World", run this command:

`$ nestor new "Hello World"`

Quotes are not needed if the name of your power is a single word.

This will create a NodeJS module with the following contents:

1. `index.js`: Contains a sample power implementation
2. `package.json`: Where you can add all your dependencies. The
   [nestorbot](https://github.com/zerobotlabs/nestorbot) dependency is
added by default.
3. `nestor.json`: A manifest file containing details about your power
4. `README.md`: Contains the [Nestor programming
   manual](https://github.com/zerobotlabs/nestorbot)

## Make Changes to your Power

You can now follow the [Nestor Programming
Manual](https://github.com/zerobotlabs/nestorbot) to make changes to
your power.

#### nestor.json

The `nestor.json` file located at the root of your power's directory contains important information that is required by Nestor:

1. `name`: The name of your power
2. `permalink`: The permalink that will uniquely identify your power.
   The permalink must not contain any spaces.
3. `Description`: A short description of what users can do with this power.

These three fields are mandatory.

In addition if your power requires environment variables that need to be configured (for e.g. authentication tokens or keys), you can set them with the `environment_keys` fields.

An example of this setting can be found in the [Mixpanel
power](https://github.com/nestor-powers/mixpanel/blob/20b779959d07a497e125b87c937011c4828f80d8/nestor.json)

```json
"environment_keys": {
  "NESTOR_MIXPANEL_API_KEY": {
    "required": true,
    "mode": "user"
  },
  "NESTOR_MIXPANEL_API_SECRET": {
    "required": true,
    "mode": "user"
  }
}
```

In this example, `NESTOR_MIXPANEL_API_KEY` and
`NESTOR_MIXPANEL_API_SECRET` are both required (specified by
`"required":true`) by the power to work
and need to be set by the user (specified by `"mode": "user"`)

By setting an environment variable to be "required", every time a user
tries to use your power, she will be prompted to set this
environment variable. This way you don't have to write additional code
in your power to check whether your environment variable is set.

If you have an optional environment variable, then set the
`required` field to false.

An example of an environment variable that is not required can be found
in the [Github
power](https://github.com/nestor-powers/github/blob/0766a895f4c9065bc5c5ef6df7245a9627ad1306/nestor.json).

## Save Your Power

To save your Nestor power, run `nestor save` inside the directory where
you power is created and this will upload the
code powering your power (pardon the pun) to Nestor's servers and your
power is now ready to be tested.

`$ nestor save`

Saving your power **does not** make it available to your Slack team. You
will still need to deploy your power (which we will cover later in this
README).

#### Sidenote: Logging In

All of the remaining operations (including this one) require you to be
authenticated with Nestor's service so you will be prompted to log in.
If you are logging in for the first time, you will need to [sign in to
the website](https://www.asknestor.me), go to "My Profile" by clicking
on your profile picture on the left bottom side, and setting your
password.

Keep note of your email address as that is required to log in.

## Test Your Power

The Nestor Toolkit provides you with an interactive shell which will let
you test your power before it is deployed. You can enter text commands
as they would appear in Slack, and see how your power behaves.

You can keep editing your power, `save`-ing your changes and keep
testing it with the shell until you are happy.

To start the shell, run the following command:

```
$ nestor shell
```

To quit the shell, enter the command `exit`.

## Deploy Your Power

Deploying your power means that your power will now be available to
your entire team. To deploy your power run the following command:

```
$ nestor deploy
```

This will give you a list of versions for your power (every time you
save your power, a new version is created) and you can pick the version
that you want to deploy.

If you want to deploy just the latest version, run:

```
$ nestor deploy --latest
```

## Bugs/Feedback

To report bugs and make feature requests, please [create an
issue](https://github.com/zerobotlabs/nestor-cli/issues).
