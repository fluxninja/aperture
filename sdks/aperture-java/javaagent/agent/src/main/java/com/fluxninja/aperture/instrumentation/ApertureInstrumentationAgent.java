package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.instrumentation.armeria.ArmeriaInstrumentationModule;
import com.fluxninja.aperture.instrumentation.netty.NettyInstrumentationModule;
import java.lang.instrument.Instrumentation;
import java.util.ArrayList;
import java.util.List;
import net.bytebuddy.agent.builder.AgentBuilder;

public class ApertureInstrumentationAgent {
    private static final List<InstrumentationModule> modules =
            new ArrayList<InstrumentationModule>() {
                {
                    add(new ArmeriaInstrumentationModule());
                    add(new NettyInstrumentationModule());
                }
            };

    public static void premain(String agentArgs, Instrumentation inst) {
        AgentBuilder agentBuilder =
                new AgentBuilder.Default()
                        .with(new AgentBuilder.InitializationStrategy.SelfInjection.Eager());

        for (InstrumentationModule module : modules) {
            for (TransformerInstrumentation tf : module.getTransformers()) {
                agentBuilder = agentBuilder.type(tf.getType()).transform(tf.getTransformer());
            }
        }
        agentBuilder.installOn(inst);
    }
}
