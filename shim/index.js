var _ref = require('./script/node_modules/nestorbot');
var Robot = _ref.Robot;
var TextMessage = _ref.TextMessage;
var User = _ref.User;
var NestorAdapter = _ref.NestorAdapter;
var Response = _ref.Response;
var Message = _ref.Message;

exports.handle = function(event, ctx) {
  var relaxEvent = event.__relax_event;
  var user = new User(relaxEvent.user_uid,
                    {
                      room: relaxEvent.channel_uid
                    });
  var msg = new TextMessage(user);
  var robot = new Robot(relaxEvent.team_uid, relaxEvent.relax_bot_uid, event.__debugMode);
  robot.requiredEnv = event.__nestor_required_env;
  robot.brain.mergeData(event.__nestor_brain);

  for(var envProp in event.__nestor_env) {
    process.env[envProp] = event.__nestor_env[envProp];
  }

  if(relaxEvent.im == true) {
    relaxEvent.text = "<@" + relaxEvent.relax_bot_uid + ">: " + relaxEvent.text;
  }

  message = new TextMessage(user, relaxEvent.text);
  require('script')(robot);

  robot.receive(message, function(done) {
    ctx.succeed({
      to_send: robot.toSend,
      to_suggest: (robot.toSuggest && robot.toSuggest.length > 0) ? robot.toSuggest : null,
      brain: robot.brain.data._private
    });
  });
}
